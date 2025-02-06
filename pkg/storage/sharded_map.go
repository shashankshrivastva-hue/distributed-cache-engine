package storage

import (
	"hash/fnv"
	"sync"
)

const ShardCount = 32

type ShardedStore struct {
	shards []*StoreShard
}

type StoreShard struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewShardedStore() *ShardedStore {
	store := &ShardedStore{shards: make([]*StoreShard, ShardCount)}
	for i := 0; i < ShardCount; i++ {
		store.shards[i] = &StoreShard{data: make(map[string][]byte)}
	}
	return store
}

func (s *ShardedStore) getShard(key string) *StoreShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return s.shards[uint(h.Sum32())%ShardCount]
}

func (s *ShardedStore) Set(key string, val []byte) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.data[key] = val
}

func (s *ShardedStore) Get(key string) ([]byte, bool) {
	shard := s.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	v, ok := shard.data[key]
	return v, ok
}
