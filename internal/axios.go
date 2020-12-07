package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Axios interface {
	Post(ctx context.Context, path string, body map[string]string) (*http.Response, error)
}

type axios struct {
	httpClient *http.Client
}

func NewAxios() *axios {
	return &axios{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
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

	return a.httpClient.Do(req)
}
