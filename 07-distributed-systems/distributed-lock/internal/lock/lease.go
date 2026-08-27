package lock

import (
	"context"
	"errors"
	"fmt"
)

var releaseScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	end

	return 0
`

var extendScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call(
			"PEXPIRE",
			KEYS[1],
			ARGV[2]
		)
	end

	return 0
`

func (l *Lease) Release(
	ctx context.Context,
) error {
	result, err := l.manager.client.Eval(
		ctx,
		releaseScript,
		[]string{l.key},
		l.token,
	).Int()
	if err != nil {
		return fmt.Errorf(
			"release redis lock: %w",
			err,
		)
	}

	if result == 0 {
		return ErrNotOwner
	}

	return nil
}

func (l *Lease) Extend(
	ctx context.Context,
) error {
	ttlMilliseconds := l.manager.config.TTL.Milliseconds()

	result, err := l.manager.client.Eval(
		ctx,
		extendScript,
		[]string{l.key},
		l.token,
		ttlMilliseconds,
	).Int()
	if err != nil {
		return fmt.Errorf(
			"extend redis lock: %w",
			err,
		)
	}

	if result == 0 {
		return errors.Join(
			ErrLockLost,
			ErrNotOwner,
		)
	}

	return nil
}
