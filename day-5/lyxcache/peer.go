package lyxcache

type PeerSelect interface {
	// PickPeer 根据传入的key，选择对应的节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
