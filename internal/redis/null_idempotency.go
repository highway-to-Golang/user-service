package redis

import (
	"context"
	"time"
)

type NullIdempotencyStorage struct{}

func (NullIdempotencyStorage) GetResult(context.Context, string) ([]byte, error) {
	return nil, nil
}

func (NullIdempotencyStorage) SaveResult(context.Context, string, []byte, time.Duration) error {
	return nil
}

func (NullIdempotencyStorage) AcquireLock(context.Context, string, time.Duration) (bool, error) {
	return true, nil
}

func (NullIdempotencyStorage) ReleaseLock(context.Context, string) error {
	return nil
}
