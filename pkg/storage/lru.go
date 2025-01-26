package storage

import (
	"container/list"
	"sync"
	"time"
)

type cacheItem struct {
	key       string
	value     []byte
	expiresAt time.Time
}

type LRUCache struct {
	mu         sync.RWMutex
	capacity   int
	items      map[string]*list.Element
	evictList  *list.List
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:  capacity,
		items:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

func (c *LRUCache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	if elem, ok := c.items[key]; ok {
		c.evictList.MoveToFront(elem)
		item := elem.Value.(*cacheItem)
		item.value = value
		item.expiresAt = exp
		return
	}

	if c.evictList.Len() >= c.capacity {
		c.evictOldest()
	}

	item := &cacheItem{key: key, value: value, expiresAt: exp}
	elem := c.evictList.PushFront(item)
	c.items[key] = elem
}

func (c *LRUCache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return nil, false
	}

	item := elem.Value.(*cacheItem)
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		c.removeElement(elem)
		return nil, false
	}

	c.evictList.MoveToFront(elem)
	return item.value, true
}

func (c *LRUCache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
		return true
	}
	return false
}

func (c *LRUCache) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)
	item := elem.Value.(*cacheItem)
	delete(c.items, item.key)
}

func (c *LRUCache) evictOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}
