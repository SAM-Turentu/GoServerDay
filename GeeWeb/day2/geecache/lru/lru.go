package lru

import "container/list"

/*
缓存
*/

type Cache struct {
	maxBytes  int64                         // 允许最大内存
	nbytes    int64                         // 当前已使用的内存
	ll        *list.List                    // 使用双向链表
	cache     map[string]*list.Element      // 值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 新增/更新
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 更新
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { //新增
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get 从双向链表中找到对应的节点，将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除访问最少的节点
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 取队首节点，从链表中删除
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 从字典中 c.cache 删除该节点的映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新当前使用的内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) // 调用回调函数
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
