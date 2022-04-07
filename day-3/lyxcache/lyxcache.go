package lyxcache

import (
	"fmt"
	"sync"
)

type GetFunc func(key string) ([]byte, error)

type Getter interface {
	Get(key string) ([]byte, error)
}

func (f GetFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	MainCache *Cache
	getter    Getter
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group, 0)
)

func New(maxByte int64, name string, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	g := &Group{
		name:      name,
		MainCache: &Cache{cacheBytes: maxByte},
		getter:    getter,
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	return groups[name]
}

func (g *Group) Get(key string) (*ByteView, error) {
	if key == "" {
		return &ByteView{}, fmt.Errorf("key is blank")
	}
	if value, ok := g.MainCache.Get(key); ok {
		return value, nil
	}
	return g.Local(key)
}

func (g *Group) Local(key string) (*ByteView, error) {
	val, err := g.getter.Get(key)
	if err != nil {
		return nil, err
	}
	new := ByteView{B: cloneBytes(val)}
	g.Addlocal(key, new)
	return &new, nil
}
func (g *Group) Addlocal(key string, val ByteView) {
	g.MainCache.Add(key, val)
}
