package main

import "sync"

type Storage struct {
	Storage map[string]string
	Mtx     sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		Storage: make(map[string]string),
	}
}

func (s *Storage) Set(k, v string) error {
	s.Mtx.Lock()
	defer s.Mtx.Unlock()
	s.Storage[k] = v
	return nil
}

func (s *Storage) Get(k string) (string, bool) {
	s.Mtx.Lock()
	defer s.Mtx.Unlock()
	v, ok := s.Storage[k]
	return v, ok
}
