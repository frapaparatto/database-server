package main

import (
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		setup      map[string]string
		wantStatus int
		wantBody   map[string]string
	}{
		{
			name:       "no parameters",
			method:     "GET",
			url:        "/get",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "missing or empty 'key' parameter",
			},
		},
		{
			name:       "wrong parameter name",
			method:     "GET",
			url:        "/get?key=a&foo=1",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "only the 'key' parameter is allowed",
			},
		},
		{
			name:       "repeated parameter",
			method:     "GET",
			url:        "/get?key=a&key=b",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "'key' must appear exactly once",
			},
		},
		{
			name:       "wrong parameter name",
			method:     "GET",
			url:        "/get?foo=1",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "missing or empty 'key' parameter",
			},
		},
		{
			name:       "never set",
			method:     "GET",
			url:        "/get?key=missing",
			wantStatus: http.StatusNotFound,
			wantBody: map[string]string{
				"error": "key 'missing' not found",
			},
		},
		{
			name:       "empty key",
			method:     "GET",
			url:        "/get?key=",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "missing or empty 'key' parameter",
			},
		},
		{
			name:       "parameter found",
			method:     "GET",
			url:        "/get?key=a",
			setup:      map[string]string{"a": "1"},
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"a": "1",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage := NewStorage()
			for k, v := range tc.setup {
				storage.Set(k, v)
			}

			r := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			h := handleGet(storage)
			h(w, r)

			if tc.wantStatus != w.Code {
				t.Errorf("want %d, got %d", tc.wantStatus, w.Code)
			}

			var b map[string]string
			if err := json.NewDecoder(w.Body).Decode(&b); err != nil {
				t.Fatalf("decoding body: %v", err)
			}

			if !maps.Equal(tc.wantBody, b) {
				t.Errorf("want %v, got %v", tc.wantBody, b)
			}
		})
	}
}

//func TestSet(t *testing.T)              {}
//func TestConcurrentAccess(t *testing.T) {}
