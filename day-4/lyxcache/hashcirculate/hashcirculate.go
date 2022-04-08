package hashcirculate

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int   // 几个虚拟节点
	keys     []int // 环
	hashMap  map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string, 0),
	}
	if hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add key代表几个真实节点，根据虚拟节点计算出来hash值,然后讲hash 值放到keys里，这样一个真实的key值可以在环中对应多个点，包括虚拟节点。
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 排序，保证顺时针环
	sort.Ints(m.keys)
}

// Get 先计算真实节点的hash值，然后顺时针找到节点的的位置，因为是环，所以要%，找到位置后，拿到还上的hash值，从hashmap中获得真实节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
