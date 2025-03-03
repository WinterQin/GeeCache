package lru

import (
	"container/list"
)

// Cache 是一个LRU缓存，它不是并发安全的
type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nbytes    int64                         // 当前已使用的内存
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // 键是字符串，值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数
}

// entry 是双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// Value 用于计算值所占用的内存大小
type Value interface {
	Len() int
}

// New 创建一个新的Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add 向缓存中添加一个值
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 如果键存在，则更新对应节点的值
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 如果键不存在，则向缓存中添加该值
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 如果超过了设定的最大内存，则移除最少访问的节点
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Get 从缓存中获取一个值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 如果键存在，将其移动到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// RemoveOldest 移除最近最少访问的节点
func (c *Cache) RemoveOldest() {
	if ele := c.ll.Back(); ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len 返回缓存中的记录数
func (c *Cache) Len() int {
	return c.ll.Len()
}
