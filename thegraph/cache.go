package thegraph

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/streamingfast/dstore"
	"go.uber.org/zap/zapcore"
)

var NotFound = errors.New("not found")

type QueryCache interface {
	Key([]string) string
	Get(ctx context.Context, block uint64, key string) ([]byte, error)
	Put(ctx context.Context, block uint64, key string, cnt []byte) error

	MarshalLogObject(encoder zapcore.ObjectEncoder) error
}

type noOpCache struct {
	missCount uint64
	hitCount  uint64
}

func (n *noOpCache) Key([]string) string {
	return ""
}

func (n *noOpCache) Get(ctx context.Context, block uint64, key string) ([]byte, error) {
	n.missCount++
	return nil, NotFound
}

func (n *noOpCache) Put(ctx context.Context, block uint64, key string, cnt []byte) error {
	return nil
}

func (n *noOpCache) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddUint64("miss_count", n.missCount)
	encoder.AddUint64("hit_count", n.hitCount)
	return nil
}

type FileCache struct {
	store       dstore.Store
	contentLock sync.RWMutex
	getCount    uint64
	missCount   uint64
	hitCount    uint64
	content     map[string][]byte
}

func (f *FileCache) Get(ctx context.Context, block uint64, cacheKey string) ([]byte, error) {
	f.getCount++

	f.contentLock.RLock()
	c, found := f.content[cacheKey]
	f.contentLock.RUnlock()
	if found {
		return c, nil
	}

	fname := filename(block, cacheKey)
	found, err := f.store.FileExists(ctx, fname)
	if err != nil {
		return nil, fmt.Errorf("failed to check file %q: %w", cacheKey, err)
	}
	if !found {
		f.missCount++
		return nil, NotFound
	}
	r, err := f.store.OpenObject(ctx, fname)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", cacheKey, err)
	}
	defer r.Close()

	cnt, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file %q: %w", cacheKey, err)
	}

	f.hitCount++
	f.contentLock.Lock()
	f.content[cacheKey] = cnt
	f.contentLock.Unlock()
	return cnt, nil
}

func (f *FileCache) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	f.contentLock.RLock()
	defer f.contentLock.RUnlock()
	encoder.AddUint64("get_count", f.getCount)
	encoder.AddUint64("miss_count", f.missCount)
	encoder.AddUint64("hit_count", f.hitCount)
	encoder.AddInt("memory_cache_size", len(f.content))
	return nil
}

func (f *FileCache) Put(ctx context.Context, block uint64, cacheKey string, cnt []byte) error {
	f.contentLock.Lock()
	f.content[cacheKey] = cnt
	f.contentLock.Unlock()

	r := bytes.NewReader(cnt)
	if err := f.store.WriteObject(ctx, filename(block, cacheKey), r); err != nil {
		return fmt.Errorf("failed to save cache file %q: %w", cacheKey, err)
	}

	return nil
}

func (f *FileCache) Key(v []string) string {
	data := []byte(strings.Join(v, ""))
	return fmt.Sprintf("%x", md5.Sum(data))
}

func filename(block uint64, cacheKey string) string {
	return fmt.Sprintf("%010d/%s.json", block, cacheKey)
}
