package thegraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"time"
)

type Graph struct {
	*http.Client
	url    string
	logger *zap.Logger
}

type Option func(*Graph) *Graph

func WithLogger(logger *zap.Logger) Option {
	return func(g *Graph) *Graph {
		g.logger = logger
		return g
	}
}

func New(graphURL string, opts ...Option) *Graph {
	g := &Graph{
		Client: newClient(),
		url:    graphURL,
		logger: zap.NewNop(),
	}

	for _, opt := range opts {
		g = opt(g)
	}
	return g
}

func (g *Graph) Fetch(ctx context.Context, query *Query) ([]byte, error) {
	g.logger.Debug("hitting thegraph api",
		zap.String("url", g.url),
		zap.String("query", query.Query),
	)

	cnt, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("unale to marshall payload: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, g.url, bytes.NewBuffer(cnt))
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := g.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %w", err)
	}
	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("querying graphql: got status %d, body %s", resp.StatusCode, string(responseBytes))
	}

	g.logger.Debug("response",
		zap.String("response", string(responseBytes)),
	)

	// TODO: this is not he best way should move to jsonb and plunk out the value
	errResp := ErrorResponse{}
	json.Unmarshal(responseBytes, &errResp)
	if len(errResp.Errors) > 0 {
		return nil, fmt.Errorf("received graphq error: %s", errResp.Errors[0].Message)
	}

	return responseBytes, nil

}

type ErrorResponse struct {
	Errors []struct {
		Locations []interface{} `json:"locations"`
		Message   string        `json:"message"`
	} `json:"errors"`
}

func newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy:              http.ProxyFromEnvironment,
			DisableKeepAlives:  false,
			DisableCompression: false,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 300 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
