package geecache

// ByteView 保存不可变的字节视图
type ByteView struct {
	b []byte
}

// Len 返回视图的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回一个拷贝，防止缓存值被外部修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 将数据作为字符串返回
func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes 返回一个byte切片的拷贝
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
