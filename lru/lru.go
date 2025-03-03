package lru

import "container/list"

// New 方便实例化 Cache，实现 New() 函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	c := &Cache{
		maxBytes:  maxBytes,
		usedBytes: 0,
		entryList: list.New(),
		kpMap:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
	return c
}

// Get 从缓存中获取一个数据
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 从kpMap中根据key，取到链表里面的指针
	LinkedListPoint, ok := c.kpMap[key]
	//如果有这个key
	if ok {
		//将这个指针移动到队头
		c.entryList.MoveToFront(LinkedListPoint)
		//指针对应的是entry，包含key 和value，利于我们从map中删除对应的key
		kvEntry := LinkedListPoint.Value.(*entry)
		//返回value，entry.value
		return kvEntry.value, true
	}
	//没有就直接返回
	return
}

// Add 向缓存中添加一个数据
func (c *Cache) Add(key string, value Value) {
	// 从kpMap中根据key，取到链表里面的指针
	LinkedListPoint, ok := c.kpMap[key]
	// 如果有这个key，那么就将这个数据移动到队头
	if ok {
		c.entryList.MoveToFront(LinkedListPoint)
		// 对数据做一个更新，包括总长度的更新，和该条数据的更新
		kvEntry := LinkedListPoint.Value.(*entry)
		// 新增的数据大小为，新增数据的长度减去原来数据的长度
		c.usedBytes += int64(value.Len()) - int64(kvEntry.value.Len())
		//更新数据
		kvEntry.value = value
		//
	} else { //如果没有这个key，就添加entry数据并且移动到队头,并且使用NewPoint接收指针
		NewPoint := c.entryList.PushFront(&entry{key, value})
		//指针添加到kpMap中
		c.kpMap[key] = NewPoint
		// 更新数据总大小
		c.usedBytes += int64(value.Len())
	}
	//如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

// RemoveOldest 移除缓存
func (c *Cache) RemoveOldest() {
	//找到队尾指针
	LinkedListPoint := c.entryList.Back()
	//不为空则开始移除
	if LinkedListPoint != nil {
		//根据entry的key，从kpMap中删除该记录
		kvEntry := LinkedListPoint.Value.(*entry)
		delete(c.kpMap, kvEntry.key)
		//直接从队列中移除
		c.entryList.Remove(LinkedListPoint)
		//维护usedBytes大小
		c.usedBytes -= int64(kvEntry.value.Len())
		//如果回调函数不为空，则运行
		if c.onEvicted != nil {
			c.onEvicted(kvEntry.key, kvEntry.value)
		}
	}
}

// Len the number of kpMap entries
func (c *Cache) Len() int {
	return c.entryList.Len()
}
