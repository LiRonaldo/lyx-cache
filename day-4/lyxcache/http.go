package lyxcache

import (
	"log"
	"net/http"
	"strings"
)

const defaultPath = "/lyxcache"

type HttpPool struct {
	self     string
	basePath string
}

func NewHttp(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultPath,
	}
}

func (h *HttpPool) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !strings.HasPrefix(request.URL.Path, h.basePath) {
		panic("request url is err")
	}
	log.Println(request.URL.Path[len(h.basePath):])
	//127.0.0.1:9999/lyxcache/source/1
	parts := strings.Split(request.URL.Path[len(h.basePath):], "/")
	log.Print(parts)
	name := parts[1]
	key := parts[2]
	group := GetGroup(name)
	if group == nil {
		panic("name is blank")
		return
	}
	v, err := group.Get(key)
	if err != nil {
		http.Error(writer, "value is blank", http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Write(v.ByteSlice())
}
