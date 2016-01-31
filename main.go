package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// If we get an INSECURE_PORT env var we forward all requests on that port to HTTPs. This allows us to get
	// around routing systems that don't forward the x-forwarded-proto headers correctly (like https://convox.com/)
	if os.Getenv("INSECURE_PORT") != "" {
		fmt.Println("Serving HTTPs redirector on port " + os.Getenv("INSECURE_PORT"))
		go http.ListenAndServe(":"+os.Getenv("INSECURE_PORT"), http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			http.Redirect(res, req, "https://"+req.Host+req.URL.String(), 301)
		}))
	}

	assets := http.StripPrefix("/", http.FileServer(http.Dir("assets/")))
	http.Handle("/", assets)
	fmt.Println("Serving on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
