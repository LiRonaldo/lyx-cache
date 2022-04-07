package lyxcache

type ByteView struct {
	B []byte
}

func (b *ByteView) Len() int {
	return len(b.B)
}

func (b *ByteView) String() string {
	return string(b.B)
}
func (b *ByteView) ByteSlice() []byte {
	return cloneBytes(b.B)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
