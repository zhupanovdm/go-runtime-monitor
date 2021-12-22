package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.HandlerFunc(metricUpdateHandler))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func metricUpdateHandler(_ http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s\n", req.Method, req.URL.Path)
}
