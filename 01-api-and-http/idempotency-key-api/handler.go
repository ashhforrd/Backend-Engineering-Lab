package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /payments", s.createPayment)

	return mux
}

func (s *Server) createPayment(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		http.Error(w, "missing Idempotency-Key", http.StatusBadRequest)
		return
	}

	var request CreatePaymentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 || request.Currency == "" {
		http.Error(w, "amount and currency are required", http.StatusBadRequest)
		return
	}

	requestHash := hashRequest(request)

	s.mu.Lock()

	if record, exist := s.records[key]; exist {
		s.mu.Unlock()

		if record.RequestHash != requestHash {
			http.Error(
				w,
				"idempotency key reused with different request",
				http.StatusConflict,
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Idempotency-Replayed", "true")
		w.WriteHeader(record.StatusCode)
		_, _ = w.Write(record.Response)
		return
	}

	s.nextPaymentID++

	payment := Payment{
		ID:       fmt.Sprintf("pay_%d", s.nextPaymentID),
		Amount:   request.Amount,
		Currency: request.Currency,
	}

	response, err := json.Marshal(payment)
	if err != nil {
		s.mu.Unlock()
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	s.records[key] = IdempotencyRecord{
		RequestHash: requestHash,
		StatusCode:  http.StatusCreated,
		Response:    response,
	}

	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}
