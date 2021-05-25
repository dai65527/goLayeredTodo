package main

import "net/http"

func main() {
	// server config
	addr := "0.0.0.0:8000"
	server := &http.Server{
		Addr: addr,
	}

	// add handlers
	http.HandleFunc("/", indexHandler)

	// start server
	server.ListenAndServe()
}
