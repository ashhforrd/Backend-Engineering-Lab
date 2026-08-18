package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type IdempotencyRecord struct {
	RequestHash string
	StatusCode  int
	Response    []byte
}

func hashRequest(req CreatePaymentRequest) string {
	data, _ := json.Marshal(req)
	sum := sha256.Sum256(data)

	return hex.EncodeToString(sum[:])
}
