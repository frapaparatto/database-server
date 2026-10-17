package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "4000"
	mux := http.NewServeMux()
	storage := NewStorage()
	s := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	getHandler := handleGet(storage)
	setHandler := handleSet(storage)

	mux.HandleFunc("GET /get", getHandler)
	mux.HandleFunc("/set", setHandler)

	log.Println("Starting database server")
	log.Fatal(s.ListenAndServe())
}
