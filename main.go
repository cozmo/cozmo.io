package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func isHttps(r *http.Request) bool {
	if r.URL.Scheme == "https" {
		return true
	}
	if strings.HasPrefix(r.Proto, "HTTPS") {
		return true
	}
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}
	return false
}

func ensureHttps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !isHttps(req) {
			hostToSend := req.Host
			if req.Header.Get("X-Forwarded-Host") != "" {
				hostToSend = req.Header.Get("X-Forwarded-Host")
			}
			http.Redirect(res, req, "https://"+hostToSend+req.URL.String(), 301)
		} else {
			next.ServeHTTP(res, req)
		}
	})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	assets := http.StripPrefix("/", http.FileServer(http.Dir("assets/")))
	if os.Getenv("ENFORCE_HTTPS") != "" {
		http.Handle("/", ensureHttps(assets))
	} else {
		http.Handle("/", assets)
	}
	fmt.Println("Serving on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
