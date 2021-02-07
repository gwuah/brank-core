package axios

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Axios struct {
	httpClient  http.Client
	bearerToken string
}

func New(client http.Client) *Axios {
	return &Axios{httpClient: client}
}

func (a *Axios) SetBearerToken(bearerToken string) {
	a.bearerToken = bearerToken
}

func (a *Axios) Post(ctx context.Context, path string, body map[string]string) (*http.Response, error) {
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

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	return a.httpClient.Do(req)
}

func (a *Axios) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if a.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+a.bearerToken)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	return a.httpClient.Do(req)
}
