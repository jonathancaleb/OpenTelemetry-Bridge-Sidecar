package main

import (
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world this is golang!"))
}

func main() {
	http.HandleFunc("/", HelloWorld)

	port := ":8080"
	http.ListenAndServe(port, nil)
}