package main

import (
	"fmt"
	"testing"
)

type Stirng string

func (s Stirng) Len() int {
	return len(s)
}

func TestCache_Add(t *testing.T) {
	lru := New(int64(4), nil)
	lru.Add("a", Stirng("1"))
	lru.Add("b", Stirng("2"))
	lru.Add("c", Stirng("3"))
	fmt.Println(lru.nBytes)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("a", Stirng("1"))
	lru.Add("b", Stirng("2"))
	lru.Add("c", Stirng("3"))
	fmt.Println(lru.Get("c"))
	fmt.Println(lru.Get("a"))
}

func TestCache_RemoveOldest(t *testing.T) {
	lru := New(int64(4), func(key string, val Value) {
		fmt.Printf("delete key:%s,value:%s", key, val)
	})
	lru.Add("a", Stirng("1"))
	lru.Add("b", Stirng("2"))
	lru.Add("c", Stirng("3"))

}
