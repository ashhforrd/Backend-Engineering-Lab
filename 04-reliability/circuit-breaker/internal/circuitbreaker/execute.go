package circuitbreaker

import "context"

type Operation[T any] func(
	ctx context.Context,
) (T, error)

type IsFailure func(err error) bool

func Execute[T any](
	ctx context.Context,
	breaker *Breaker,
	operation Operation[T],
	isFailure IsFailure,
) (T, error) {
	var zero T

	requestTicket, err := breaker.beforeRequest()
	if err != nil {
		return zero, err
	}

	result, err := operation(ctx)
	if err == nil {
		breaker.afterSuccess(requestTicket)
		return result, nil
	}

	if isFailure(err) {
		breaker.afterFailure(requestTicket)
	} else {
		breaker.afterSuccess(requestTicket)
	}

	return zero, err
}