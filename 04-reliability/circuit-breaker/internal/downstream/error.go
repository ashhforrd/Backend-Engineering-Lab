package downstream

import (
	"errors"
	"fmt"
)

type Error struct {
	StatusCode int
	Message string
	Failure bool
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"downstream failed: status=%d messsage=%s",
		e.StatusCode,
		e.Message,
	)
}

func IsFailure(err error) bool {
	var downstreamErr *Error

	if !errors.As(err, &downstreamErr) {
		return false
	}

	return downstreamErr.Failure
}