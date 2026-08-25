package downstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL string
	httpClient *http.Client
}

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Price int64 `json:"price"`
}

func NewClient(
	baseURL string,
	httpClient *http.Client,
) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *Client) GetProduct(
	ctx context.Context,
	productID string,
) (Product, error) {
	var product Product

	requestURL := c.baseURL + 
		"/products/" +
		url.PathEscape(productID)
	
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL,
		nil,
	)
	if err != nil {
		return product, fmt.Errorf(
			"create downstream request: %w",
			err,
		)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return product, &Error{
			StatusCode: 0,
			Message: err.Error(),
			Failure: true,
		}
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || 
		response.StatusCode >= 300 {
			message, readErr := io.ReadAll(
				io.LimitReader(response.Body, 4*1024),
			)
			if readErr != nil {
				message = []byte("unable to read downstream error")
			}

			return product, &Error{
				StatusCode: response.StatusCode,
				Message: string(message),
				Failure: isFailureStatus(response.StatusCode),
			}
		}

	if err := json.NewDecoder(response.Body).Decode(&product); err != nil {
		return product, &Error{
			StatusCode: response.StatusCode,
			Message: "invalid downstream response",
			Failure: true,
		}
	}

	return product, nil
}

func isFailureStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}