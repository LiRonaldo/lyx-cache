package lyxcache

import (
	"day-5/lyxcache/hashcirculate"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/lyxcache"
	defaultReplicas = 50
)

// HttpServer 服务端
type HttpServer struct {
	self     string
	basePath string
	mu       sync.Mutex
	// hash 环
	perrs *hashcirculate.Map
	node  map[string]*HttpClient
}

func NewHttp(self string) *HttpServer {
	return &HttpServer{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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

// HttpClient 客户端
type HttpClient struct {
	baseURL string
}

var _ PeerGetter = (*HttpClient)(nil)

func (h *HttpClient) Get(group string, key string) ([]byte, error) {
	reqUrl := fmt.Sprintf("%v%v%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return %s", err)
	}
	b, err := ioutil.ReadAll(res.Body)
	return b, err

}

// Set 设置hash环上的节点
func (h *HttpServer) Set(peer ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.perrs = hashcirculate.New(defaultReplicas, nil)
	h.perrs.Add(peer...)
	h.node = make(map[string]*HttpClient, len(peer))
	for _, v := range peer {
		h.node[v] = &HttpClient{baseURL: h.basePath + v}
	}
}

// PickPeer 选择节点
func (h *HttpServer) PickPeer(key string) (peer PeerGetter, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.perrs.Get(key); peer != "" {
		// 因为传入的key是真实节点，通过get()获得真实的节点，
		//为啥不直接h.node[peer]获得，因为环上有虚拟节点
		return h.node[peer], true
	}
	return nil, false
}
