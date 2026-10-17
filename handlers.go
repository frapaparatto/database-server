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
		k := d.Get("key")

		if k == "" {
			writeErrorResponse(w, http.StatusBadRequest, "missing or empty 'key' parameter", nil)
			return
		}

		if len(d) != 1 {
			writeErrorResponse(w, http.StatusBadRequest, "only the 'key' parameter is allowed", nil)
			return
		}

		if len(d["key"]) != 1 {
			writeErrorResponse(w, http.StatusBadRequest, "'key' must appear exactly once", nil)
			return
		}

		v, ok := s.Get(k)
		if !ok {
			writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("key '%s' not found", k), nil)
			return
		}
		writeJSONResponse(w, http.StatusOK, map[string]string{k: v})
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
