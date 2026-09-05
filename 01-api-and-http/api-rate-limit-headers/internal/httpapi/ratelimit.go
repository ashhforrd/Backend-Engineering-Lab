package httpapi

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-rate-limit-headers/internal/ratelimit"
)

func RateLimitMiddleware(
	limiter *ratelimit.Limiter,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			clientID := identifyClient(request)
			decision := limiter.Allow(clientID)

			writer.Header().Set(
				"X-RateLimit-Limit",
				strconv.Itoa(decision.Limit),
			)
			writer.Header().Set(
				"X-RateLimit-Remaining",
				strconv.Itoa(decision.Remaining),
			)
			writer.Header().Set(
				"X-RateLimit-Reset",
				strconv.FormatInt(
					decision.ResetAt.Unix(),
					10,
				),
			)

			if !decision.Allowed {
				retryAfterSeconds := int64(
					(decision.RetryAfter +
						time.Second - 1) /
						time.Second,
				)

				if retryAfterSeconds < 1 {
					retryAfterSeconds = 1
				}

				writer.Header().Set(
					"Retry-After",
					strconv.FormatInt(
						retryAfterSeconds,
						10,
					),
				)

				http.Error(
					writer,
					"rate limit exceeded",
					http.StatusTooManyRequests,
				)

				return
			}

			next.ServeHTTP(writer, request)
		},
	)
}

func identifyClient(
	request *http.Request,
) string {
	if apiKey := request.Header.Get("X-API-Key"); apiKey != "" {
		return "api-key:" + apiKey
	}

	host, _, err := net.SplitHostPort(
		request.RemoteAddr,
	)
	if err != nil {
		return "ip:" + request.RemoteAddr
	}

	return "ip:" + host
}
