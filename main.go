package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	assets := http.StripPrefix("/", http.FileServer(http.Dir("assets/")))
	http.Handle("/", assets)
	fmt.Println("Serving on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
