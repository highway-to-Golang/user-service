package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/highway-to-Golang/user-service/internal/domain"
	apperrors "github.com/highway-to-Golang/user-service/internal/errors"
)

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *HTTPClient) CreateUser(ctx context.Context, req domain.CreateUserRequest, idempotencyKey string) (domain.User, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/users", bytes.NewReader(body))
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if idempotencyKey != "" {
		httpReq.Header.Set("Idempotency-Key", idempotencyKey)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return domain.User{}, c.mapHTTPError(resp.StatusCode, bodyBytes)
	}

	var user domain.User
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return domain.User{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return user, nil
}

func (c *HTTPClient) GetUser(ctx context.Context, id string) (domain.User, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/users/"+id, nil)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return domain.User{}, c.mapHTTPError(resp.StatusCode, bodyBytes)
	}

	var user domain.User
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return domain.User{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return user, nil
}

func (c *HTTPClient) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/users", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.mapHTTPError(resp.StatusCode, bodyBytes)
	}

	var response struct {
		Users []domain.User `json:"users"`
		Total int           `json:"total"`
	}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Users, nil
}

func (c *HTTPClient) UpdateUser(ctx context.Context, id string, req domain.UpdateUserRequest) (domain.User, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/api/users/"+id, bytes.NewReader(body))
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return domain.User{}, c.mapHTTPError(resp.StatusCode, bodyBytes)
	}

	var user domain.User
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return domain.User{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return user, nil
}

func (c *HTTPClient) DeleteUser(ctx context.Context, id string) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/api/users/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return c.mapHTTPError(resp.StatusCode, bodyBytes)
	}

	return nil
}

func (c *HTTPClient) Close() error {
	c.client.CloseIdleConnections()
	return nil
}

func (c *HTTPClient) mapHTTPError(statusCode int, body []byte) error {
	var errorResp struct {
		Error string `json:"error"`
	}
	_ = json.Unmarshal(body, &errorResp)

	switch statusCode {
	case http.StatusNotFound:
		return domain.ErrNotFound
	case http.StatusUnprocessableEntity:
		if errorResp.Error == "Request already in progress" {
			return apperrors.ErrRequestAlreadyInProgress
		}
		return fmt.Errorf("unprocessable entity: %s", errorResp.Error)
	default:
		return fmt.Errorf("HTTP error %d: %s", statusCode, errorResp.Error)
	}
}
