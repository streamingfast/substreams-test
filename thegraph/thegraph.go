package thegraph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/streamingfast/dstore"
	"go.uber.org/zap"
)

type Client struct {
	*http.Client
	url   string
	cache QueryCache

	cacheHitCount  uint64
	cacheMissCount uint64
	zlogger        *zap.Logger
	count          int
	retries        int
}

type Option func(*Client) *Client

func WithCache(cacheStore dstore.Store) Option {
	return func(g *Client) *Client {
		g.cache = &FileCache{
			store:   cacheStore,
			content: map[string][]byte{},
		}
		return g
	}
}

func WithLogger(zlogger *zap.Logger) Option {
	return func(g *Client) *Client {
		g.zlogger = zlogger
		return g
	}
}

func New(graphURL string, opts ...Option) *Client {
	g := &Client{
		Client:  newClient(),
		url:     graphURL,
		cache:   &noOpCache{},
		zlogger: zap.NewNop(),
		count:   0,
		retries: 5,
	}

	for _, opt := range opts {
		g = opt(g)
	}
	return g
}

func (g *Client) Fetch(ctx context.Context, blockNum uint64, query string, vars map[string]interface{}) ([]byte, error) {

	chunk := []string{
		g.url,
	}
	for k, v := range vars {
		chunk = append(chunk, fmt.Sprintf("%s=%s", k, v))
	}
	cacheKey := g.cache.Key(chunk)

	cnt, err := g.cache.Get(ctx, blockNum, cacheKey)
	if err == nil {
		g.zlogger.Debug("cache hit", zap.String("cache_key", cacheKey))
		g.cacheHitCount++
		return cnt, nil
	}

	g.zlogger.Debug("cache misses", zap.String("cache_key", cacheKey), zap.String("error", err.Error()))

	params := map[string]interface{}{
		"query":     query,
		"variables": vars,
	}
	g.cacheMissCount++
	g.zlogger.Debug("http fetching", zap.Reflect("params", params))
	cnt, err = g.fetch(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch entities: %w", err)
	}
	if err := g.cache.Put(ctx, blockNum, cacheKey, cnt); err != nil {
		g.zlogger.Warn("cache put failed", zap.Error(err))
	}
	return cnt, nil
}

func (g *Client) fetch(ctx context.Context, payload map[string]interface{}) ([]byte, error) {
	g.zlogger.Debug("hitting thegraph api",
		zap.String("url", g.url),
		zap.String("query", payload["query"].(string)),
	)

	cnt, err := json.Marshal(payload)
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

	g.zlogger.Debug("response",
		zap.String("response", string(responseBytes)),
	)

	// TODO: this is not he best way should move to jsonb and plunk out the value
	errResp := ErrorResponse{}
	err = json.Unmarshal(responseBytes, &errResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	for i := 0; i < g.retries; i++ {
		if len(errResp.Errors) > 0 {
			g.zlogger.Info(fmt.Sprintf("retry number: %d waiting 3s to not overload the graphAPI", i))
			time.Sleep(3 * time.Second)

			responseBytes, err = g.retryCall(request)
			if err != nil {
				return nil, fmt.Errorf("retry call: %w", err)
			}

			errResp = ErrorResponse{}
			err = json.Unmarshal(responseBytes, &errResp)
			if err != nil {
				return nil, fmt.Errorf("unmarshalling: %w", err)
			}
		}
	}

	return responseBytes, nil
}

func (g *Client) retryCall(request *http.Request) ([]byte, error) {
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

func (g *Client) GetCacheHitCount() uint64 {
	return g.cacheHitCount
}

func (g *Client) GetCacheMissCount() uint64 {
	return g.cacheMissCount
}

func (g *Client) GetTotalQueries() uint64 {
	return g.cacheHitCount + g.cacheMissCount
}
