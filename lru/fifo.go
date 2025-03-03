package lru

import "container/list"

func NewFifoCache(maxSize int64) *fifoCache {
	return &fifoCache{
		maxSize:   maxSize,
		usedBytes: 0,
		queue:     list.New(),
		kpMap:     make(map[string]*list.Element),
		onEvicted: nil,
	}
}

// Len the number of kpMap entries
func (fc *fifoCache) Len() int {
	return fc.queue.Len()
}

func (fc *fifoCache) Get(key string) (value Value, ok bool) {
	// 先查指针
	point, ok := fc.kpMap[key]
	//如果找到了
	if ok {
		kvEntry := point.Value.(*entry)
		return kvEntry.value, ok
	}
	//没找到就直接返回
	return
}

func (fc *fifoCache) Add(key string, value Value) {
	// 先查指针
	point, ok := fc.kpMap[key]
	//如果找到了,就修改值
	if ok {
		//通过指针找到entry
		kvEntry := point.Value.(*entry)
		//修改已使用大小，用新增的大小减去原来的大小
		fc.usedBytes += int64(value.Len() - kvEntry.value.Len())
		//修改值
		kvEntry.value = value

	} else { //如果没找到就增加到队头
		newPoint := fc.queue.PushFront(&entry{key: key, value: value})
		//将其指针存放到kpMap里面
		fc.kpMap[key] = newPoint
		// 更新已使用空间
		fc.usedBytes += int64(value.Len())
	}
	for fc.maxSize > 0 && fc.usedBytes >= fc.maxSize {
		fc.RemoveOldest()
	}
}

func (fc *fifoCache) RemoveOldest() {
	point := fc.queue.Back()
	if point != nil {
		//根据entry的key，从kpMap中删除该记录
		kvEntry := point.Value.(*entry)
		delete(fc.kpMap, kvEntry.key)
		//直接从队列中移除
		fc.queue.Remove(point)
		//维护usedBytes大小
		fc.usedBytes -= int64(kvEntry.value.Len())
		//如果回调函数不为空，则运行
		if fc.onEvicted != nil {
			fc.onEvicted(kvEntry.key, kvEntry.value)
		}
	}
}
