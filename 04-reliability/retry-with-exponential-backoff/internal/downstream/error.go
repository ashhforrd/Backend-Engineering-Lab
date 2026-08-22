package downstream

import "fmt"

type Error struct {
	StatusCode int
	Message    string
	Retryable  bool
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"downstream requesst failed: statu: %d message: %s",
		e.StatusCode,
		e.Message,
	)
}

func IsRetryable(err error) bool {
	downstreamErr, ok := err.(*Error)
	if !ok {
		return false
	}

	return downstreamErr.Retryable
}
