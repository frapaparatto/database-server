package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleGet(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("URL of the request: ", r.URL)
		d := r.URL.Query()

		if v, found := s.Get(d.Get("key")); found {
			io.WriteString(w, fmt.Sprintf("key: %s, value: %s\n", d.Get("key"), v))
		} else {
			io.WriteString(w, fmt.Sprintf("key: %s not found", d.Get("key")))
		}
	}
}

func handleSet(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("URL of the request: ", r.URL)
		q := r.URL.Query()
		for k, v := range q {
			s.Set(k, v[0])
		}

		io.WriteString(w, fmt.Sprintf("Storage is: %v\n", s))
	}
}
