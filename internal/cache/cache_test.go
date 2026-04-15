package cache_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/cache"
	"github.com/driftwatch/internal/drift"
)

func makeResults(image string) []drift.Result {
	return []drift.Result{
		{Service: "web", ExpectedImage: image, ActualImage: image, Drifted: false},
	}
}

func TestGet_MissOnEmptyCache(t *testing.T) {
	c := cache.New(5 * time.Second)
	_, ok := c.Get("web")
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSet_And_Get_Hit(t *testing.T) {
	c := cache.New(5 * time.Second)
	results := makeResults("nginx:latest")
	c.Set("web", results)

	got, ok := c.Get("web")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != 1 || got[0].ExpectedImage != "nginx:latest" {
		t.Errorf("unexpected results: %+v", got)
	}
}

func TestGet_Miss_AfterExpiry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("web", makeResults("nginx:latest"))
	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("web")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestGet_Miss_ZeroTTL(t *testing.T) {
	c := cache.New(0)
	c.Set("web", makeResults("nginx:latest"))

	_, ok := c.Get("web")
	if ok {
		t.Fatal("expected cache miss when TTL is zero")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("web", makeResults("nginx:latest"))
	c.Invalidate("web")

	_, ok := c.Get("web")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(5 * time.Second)
	c.Set("web", makeResults("nginx:latest"))
	c.Set("api", makeResults("alpine:3.18"))
	c.Flush()

	for _, key := range []string{"web", "api"} {
		if _, ok := c.Get(key); ok {
			t.Errorf("expected cache miss for %q after flush", key)
		}
	}
}
