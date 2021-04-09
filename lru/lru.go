package lru

import "container/list"

type Cache struct {
	maxBytes int64
	ll       *list.List
	cache    map[string]*list.Element
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		cache := ele.Value.(*entry)
		delete(c.cache, cache.key)
	}
}

func (c *Cache) Set(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		cache := ele.Value.(*entry)
		cache.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
	}
	for c.maxBytes != 0 && int64(c.Len()) > c.maxBytes {
		delete(c.cache, c.ll.Remove(c.ll.Back()).(*entry).key)
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
