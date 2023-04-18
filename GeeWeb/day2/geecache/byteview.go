package geecache

type ByteView struct {
	b []byte //存储真实的缓存值，byte 可以支持任意类型存储
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b) // 返回拷贝值，可以防止缓存值被外部程序修改
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
