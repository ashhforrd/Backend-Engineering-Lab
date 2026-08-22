package downstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type PaymentRequest struct {
	OrderID               string `json:"orderId"`
	Amount                int64  `json:"amount"`
	FailuresBeforeSuccess int    `json:"failuresBeforeSuccess"`
}

type PaymentResponse struct {
	PaymentID string `json:"paymentId"`
	Status    string `json:"status"`
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *Client) CreatePayment(
	ctx context.Context,
	request PaymentRequest,
) (PaymentResponse, error) {
	var response PaymentResponse

	body, err := json.Marshal(request)
	if err != nil {
		return response, fmt.Errorf("encode payment request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/payments",
		bytes.NewReader(body),
	)
	if err != nil {
		return response, fmt.Errorf("create payment request: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return response, &Error{
			StatusCode: 0,
			Message:    err.Error(),
			Retryable:  true,
		}
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		message, readErr := io.ReadAll(
			io.LimitReader(httpResponse.Body, 4*1024),
		)
		if readErr != nil {
			message = []byte("unable to read downstream error")
		}

		return response, &Error{
			StatusCode: httpResponse.StatusCode,
			Message:    string(message),
			Retryable:  isRetryableStatus(httpResponse.StatusCode),
		}
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("decode payment response: %w", err)
	}

	return response, nil
}

func isRetryableStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
