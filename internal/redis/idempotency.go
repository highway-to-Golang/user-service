package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type IdempotencyStorage struct {
	client *redis.Client
}

func NewIdempotencyStorage(url string) (*IdempotencyStorage, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &IdempotencyStorage{client: client}, nil
}

func (s *IdempotencyStorage) Close() error {
	return s.client.Close()
}

func (s *IdempotencyStorage) GetResult(ctx context.Context, key string) ([]byte, error) {
	data, err := s.client.Get(ctx, fmt.Sprintf("idempotency:%s", key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	return data, err
}

func (s *IdempotencyStorage) SaveResult(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return s.client.Set(ctx, fmt.Sprintf("idempotency:%s", key), value, ttl).Err()
}

func (s *IdempotencyStorage) AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	result := s.client.SetNX(ctx, fmt.Sprintf("lock:%s", key), "1", ttl)
	return result.Val(), result.Err()
}

func (s *IdempotencyStorage) ReleaseLock(ctx context.Context, key string) error {
	return s.client.Del(ctx, fmt.Sprintf("lock:%s", key)).Err()
}
