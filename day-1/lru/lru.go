package main

import "container/list"

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entity struct {
	key   string
	value Value
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element, 0),
		OnEvicted: onEvicted,
	}
}

// Add 利用转型的特性，value是接口父类，只要任何实现了Len()方法的对象，都可以作为参数,len方法是调用子类自己的
func (c *Cache) Add(key string, value Value) {
	if val, ok := c.cache[key]; ok {
		c.ll.MoveToFront(val)
		kv := val.Value.(*entity)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		val := c.ll.PushFront(&entity{key: key, value: value})
		c.cache[key] = val
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	if c.maxBytes != 0 && c.nBytes > c.maxBytes {

	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if val, ok := c.cache[key]; ok {
		c.ll.MoveToFront(val)
		e := val.Value.(*entity)
		return e.value, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	e := c.ll.Back()
	if e != nil {
		kv := e.Value.(*entity)
		delete(c.cache, kv.key)
		c.ll.Remove(e)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

type Value interface {
	Len() int
}
