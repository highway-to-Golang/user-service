package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

type EventSink struct {
	conn          *nats.Conn
	subjectPrefix string
}

func New(url, subjectPrefix string) (*EventSink, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &EventSink{
		conn:          conn,
		subjectPrefix: subjectPrefix,
	}, nil
}

func (es *EventSink) Close() {
	if es.conn != nil {
		es.conn.Close()
	}
}

type Event struct {
	Method    string    `json:"method"`
	Timestamp time.Time `json:"timestamp"`
}

func (es *EventSink) Publish(ctx context.Context, method string) error {
	event := Event{
		Method:    method,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	subject := fmt.Sprintf("%s.%s", es.subjectPrefix, method)
	if err := es.conn.Publish(subject, data); err != nil {
		slog.Error("failed to publish event", "error", err, "subject", subject)
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
