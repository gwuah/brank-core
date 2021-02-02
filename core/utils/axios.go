package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Axios interface {
	Get(ctx context.Context, path string) (*http.Response, error)
	Post(ctx context.Context, path string, body map[string]string) (*http.Response, error)
	SetBearerToken(bearerToken string)
}

type axios struct {
	httpClient  *http.Client
	bearerToken string
}

func NewAxios() *axios {
	return &axios{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *axios) SetBearerToken(bearerToken string) {
	a.bearerToken = bearerToken
}

func (a *axios) Post(ctx context.Context, path string, body map[string]string) (*http.Response, error) {
	rb, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	if a.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+a.bearerToken)
	}

	return a.httpClient.Do(req)
}

func (a *axios) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if a.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+a.bearerToken)
	}

	return a.httpClient.Do(req)
}
