package lru

import "container/list"

// 键值对 entry 是双向链表节点的数据类型，
// 在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
type entry struct {
	key   string
	value Value
}

// Value 为了通用性，我们允许值是实现了 Value 接口的任意类型，
// 该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。
type Value interface {
	Len() int
}

type Cache struct {
	maxBytes  int64                    //最大存储字节
	usedBytes int64                    //当前使用的字节
	entryList *list.List               //存放entry的链表，当作队列
	kpMap     map[string]*list.Element //存放key-链表指针的map

	onEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为 nil。
}

type fifoCache struct {
	maxSize   int64
	usedBytes int64

	queue *list.List
	kpMap map[string]*list.Element

	onEvicted func(key string, value interface{})
}

type lfuCache struct {
	maxSize   int64
	usedBytes int64

	queue *list.List
	kpMap map[string]*list.Element

	onEvicted func(key string, value interface{})
}

type ICache interface {
	Get(key string) (value Value, ok bool)
	RemoveOldest()
	Add(key string, value Value)
}
