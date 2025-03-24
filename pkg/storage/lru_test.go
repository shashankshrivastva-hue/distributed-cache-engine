package storage

import (
	"testing"
	"time"
)

func TestLRUCache_Basic(t *testing.T) {
	c := NewLRUCache(2)
	c.Set("k1", []byte("v1"), 0)
	c.Set("k2", []byte("v2"), 0)

	val, ok := c.Get("k1")
	if !ok || string(val) != "v1" {
		t.Fatalf("Expected v1, got %s", val)
	}

	c.Set("k3", []byte("v3"), 0) // Should evict k2
	_, ok = c.Get("k2")
	if ok {
		t.Fatalf("Expected k2 to be evicted")
	}
}

func BenchmarkLRUCache_Get(b *testing.B) {
	c := NewLRUCache(1000)
	c.Set("key", []byte("value"), 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get("key")
	}
}
