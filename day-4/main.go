package main

import (
	"day-4/lyxcache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {

	self := "127.0.0.1:9999"
	pool := lyxcache.NewHttp(self)
	lyxcache.New(2>>10, "source", lyxcache.GetFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s is not exit", key)
	}))
	g := lyxcache.GetGroup("source")
	key := "lyx"
	g.MainCache.Add(key, lyxcache.ByteView{[]byte{1}})
	log.Fatalln(http.ListenAndServe(self, pool))

}
