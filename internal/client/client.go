package client

import (
	"context"
	"fmt"
	"github.com/gojek/heimdall/httpclient"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *httpclient.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpclient.NewClient(httpclient.WithHTTPTimeout(timeout)),
	}
}

func (c *Client) makeRequest(
	ctx context.Context,
	method, endpoint string,
	queryParams url.Values,
	headers map[string]string, // Новый параметр для заголовков
	body io.Reader,
) (*http.Response, error) {
	fullURL := c.BaseURL + endpoint
	if queryParams != nil && len(queryParams) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryParams.Encode())
	}

	request, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return response, nil
}
