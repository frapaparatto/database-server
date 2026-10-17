package main

import (
	"fmt"
	"net/http"
)

func handleGet(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		q := r.URL.Query()

		if len(q) == 0 {
			writeErrorResponse(w, http.StatusBadRequest, "no parameters provided", nil)
			return
		}

		support := make(map[string]string)
		for k, v := range q {
			if k == "" {
				writeErrorResponse(w, http.StatusBadRequest, "empty key is not allowed", nil)
				return
			}

			if v[0] == "" {
				writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("empty value for key '%s'", k), nil)
				return
			}

			if len(v) != 1 {
				writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("key '%s' must appear exactly once", k), nil)
				return
			}

			support[k] = v[0]
		}

		s.SetMany(support)
		writeJSONResponse(w, http.StatusOK, support)
	}
}
