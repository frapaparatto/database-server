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
			name:       "extra parameter",
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

			if got := w.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("want Content-Type application/json, got %q", got)
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

func TestSet(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		setup      map[string]string
		wantStatus int
		wantBody   map[string]string
		// wantStored is the full content of the store after the request.
		// nil means the store must be empty.
		wantStored map[string]string
	}{
		{
			name:       "no parameters",
			method:     "GET",
			url:        "/set",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "no parameters provided",
			},
		},
		{
			name:       "empty key",
			method:     "GET",
			url:        "/set?=1",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "empty key is not allowed",
			},
		},
		{
			name:       "empty value",
			method:     "GET",
			url:        "/set?a=",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "empty value for key 'a'",
			},
		},
		{
			name:       "key without equals sign",
			method:     "GET",
			url:        "/set?a",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "empty value for key 'a'",
			},
		},
		{
			name:       "repeated key",
			method:     "GET",
			url:        "/set?a=1&a=2",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "key 'a' must appear exactly once",
			},
		},
		{
			name:       "repeated key with second value empty",
			method:     "GET",
			url:        "/set?a=1&a=",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "key 'a' must appear exactly once",
			},
		},
		// The empty-value check runs before the repeated-key check,
		// so an empty first value is reported as an empty value.
		{
			name:       "repeated key with first value empty",
			method:     "GET",
			url:        "/set?a=&a=1",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "empty value for key 'a'",
			},
		},
		// Map iteration order is random: this row catches a handler that stores
		// pairs before validating all of them only when "a" is visited before "b".
		// It never fails on a correct handler, but it can miss the regression on some runs.
		{
			name:       "one invalid pair among valid ones",
			method:     "GET",
			url:        "/set?a=1&b=",
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "empty value for key 'b'",
			},
		},
		{
			name:       "rejected request keeps existing value",
			method:     "GET",
			url:        "/set?a=3&a=4",
			setup:      map[string]string{"a": "2"},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "key 'a' must appear exactly once",
			},
			wantStored: map[string]string{"a": "2"},
		},
		{
			name:       "single pair",
			method:     "GET",
			url:        "/set?a=1",
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"a": "1",
			},
			wantStored: map[string]string{"a": "1"},
		},
		{
			name:       "multiple pairs",
			method:     "GET",
			url:        "/set?a=1&b=2",
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"a": "1",
				"b": "2",
			},
			wantStored: map[string]string{"a": "1", "b": "2"},
		},
		{
			name:       "overwrite existing key",
			method:     "GET",
			url:        "/set?a=new",
			setup:      map[string]string{"a": "old"},
			wantStatus: http.StatusOK,
			wantBody: map[string]string{
				"a": "new",
			},
			wantStored: map[string]string{"a": "new"},
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

			h := handleSet(storage)
			h(w, r)

			if tc.wantStatus != w.Code {
				t.Errorf("want %d, got %d", tc.wantStatus, w.Code)
			}

			if got := w.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("want Content-Type application/json, got %q", got)
			}

			var b map[string]string
			if err := json.NewDecoder(w.Body).Decode(&b); err != nil {
				t.Fatalf("decoding body: %v", err)
			}

			if !maps.Equal(tc.wantBody, b) {
				t.Errorf("want %v, got %v", tc.wantBody, b)
			}

			if !maps.Equal(tc.wantStored, storage.data) {
				t.Errorf("want stored %v, got %v", tc.wantStored, storage.data)
			}
		})
	}
}

//func TestConcurrentAccess(t *testing.T) {}
