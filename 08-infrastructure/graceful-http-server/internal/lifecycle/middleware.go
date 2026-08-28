package lifecycle

import (
	"net/http"
	"sync/atomic"
)

type Middleware struct {
	state          *State
	activeRequests atomic.Int64
}

func NewMiddleware(
	state *State,
) *Middleware {
	return &Middleware{
		state: state,
	}
}

func (m *Middleware) Track(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			if !m.state.IsReady() {
				http.Error(
					writer,
					"server is shutting down",
					http.StatusServiceUnavailable,
				)
				return
			}

			m.activeRequests.Add(1)
			defer m.activeRequests.Add(-1)

			next.ServeHTTP(writer, request)
		},
	)
}

func (m *Middleware) ActiveRequests() int64 {
	return m.activeRequests.Load()
}
