package data

import (
	"container/list"
	"sync"
)

// CacheItem 缓存项
type CacheItem struct {
	Key   string
	Value interface{}
}

// LRUCache LRU缓存实现
type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	list     *list.List
	mu       sync.RWMutex
}

// NewLRUCache 创建LRU缓存
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Get 获取缓存项
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*CacheItem).Value, true
	}
	return nil, false
}

// Put 放置缓存项
func (c *LRUCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*CacheItem).Value = value
		return
	}

	// 如果缓存已满，删除最久未使用的项
	if c.list.Len() >= c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			delete(c.items, oldest.Value.(*CacheItem).Key)
			c.list.Remove(oldest)
		}
	}

	// 添加新项
	item := &CacheItem{Key: key, Value: value}
	elem := c.list.PushFront(item)
	c.items[key] = elem
}

// Delete 删除缓存项
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		delete(c.items, key)
		c.list.Remove(elem)
	}
}

// Clear 清空缓存
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.list = list.New()
}

// Len 获取缓存项数量
func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.list.Len()
}
