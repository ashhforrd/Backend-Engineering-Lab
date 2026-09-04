package httpapi

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(
	writer http.ResponseWriter,
	statusCode int,
	value any,
) {
	writer.Header().Set(
		"Content-Type",
		"application/json",
	)
	writer.WriteHeader(statusCode)

	_ = json.NewEncoder(writer).Encode(value)
}

func WriteError(
	writer http.ResponseWriter,
	statusCode int,
	message string,
) {
	WriteJSON(
		writer,
		statusCode,
		ErrorResponse{
			Error: message,
		},
	)
}
