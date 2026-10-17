package main

import "sync"

type Storage struct {
	storage map[string]string
	mtx     sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		storage: make(map[string]string),
	}
}

func (s *Storage) Set(k, v string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.storage[k] = v
	return nil
}

func (s *Storage) Get(k string) (string, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	v, ok := s.storage[k]
	return v, ok
}
