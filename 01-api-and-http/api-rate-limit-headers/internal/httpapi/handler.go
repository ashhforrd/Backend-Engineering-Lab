package httpapi

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func MessageHandler() http.Handler {
	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			_ *http.Request,
		) {
			writer.Header().Set(
				"Content-Type",
				"application/json",
			)
			writer.WriteHeader(http.StatusOK)

			_ = json.NewEncoder(writer).Encode(
				MessageResponse{
					Message: "request accepted",
				},
			)
		},
	)
}
