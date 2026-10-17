package main

import "sync"

type Storage struct {
	data map[string]string
	mtx  sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// Set has no error return: with an in-memory map it cannot fail. Add one when a failing backend (e.g. persistence to disk) is introduced.
func (s *Storage) Set(k, v string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.data[k] = v
}

func (s *Storage) Get(k string) (string, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	v, ok := s.data[k]
	return v, ok
}
