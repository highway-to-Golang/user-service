package nats

import "context"

type NullEventSink struct{}

func (NullEventSink) Publish(ctx context.Context, method string) error {
	return nil
}
