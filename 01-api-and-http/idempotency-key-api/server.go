package main

import "sync"

type Server struct {
	mu            sync.Mutex
	records       map[string]IdempotencyRecord
	nextPaymentID int64
}

func NewServer() *Server {
	return &Server{
		records: make(map[string]IdempotencyRecord),
	}
}
