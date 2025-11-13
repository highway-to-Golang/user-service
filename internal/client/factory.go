package client

import (
	"fmt"

	"github.com/highway-to-Golang/user-service/config"
)

func NewClient(cfg config.Client) (Client, error) {
	switch cfg.Protocol {
	case "http":
		return NewHTTPClient(cfg.HTTPURL), nil
	case "grpc":
		return NewGRPCClient(cfg.GRPCAddr)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s (supported: http, grpc)", cfg.Protocol)
	}
}
