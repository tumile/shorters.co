package cache

import (
	"container/list"
	"log"
)

type LFUCache interface {
	Get(key interface{}) interface{}
	Put(key, value interface{})
	Remove(key interface{}) interface{}
}

type entry struct {
	key    interface{}
	value  interface{}
	bucket int
	node   *list.Element
}

type lfuCache struct {
	freqList []list.List
	dict     map[interface{}]entry
	size     int
	capacity int
	minFreq  int
}

func NewLFUCache(capacity int) LFUCache {
	if capacity == 0 {
		log.Fatal("Capacity should not be 0")
	}
	return &lfuCache{
		freqList: make([]list.List, capacity),
		dict:     make(map[interface{}]entry),
		size:     0,
		capacity: capacity,
		minFreq:  0,
	}
}

func (c *lfuCache) Get(key interface{}) interface{} {
	if e, ok := c.dict[key]; ok {
		c.moveToNextBucket(e)
		return e.value
	}
	return nil
}

func (c *lfuCache) Put(key, value interface{}) {
	if e, ok := c.dict[key]; ok {
		e.value = value
		c.moveToNextBucket(e)
		return
	}
	if c.size == c.capacity {
		minFreqList := c.freqList[c.minFreq]
		delete(c.dict, minFreqList.Remove(minFreqList.Back()))
		c.size--
	}
	e := entry{key, value, 0, c.freqList[0].PushFront(key)}
	c.dict[key] = e
	c.minFreq = 0
	c.size++
}

func (c *lfuCache) Remove(key interface{}) interface{} {
	if e, ok := c.dict[key]; ok {
		delete(c.dict, c.freqList[e.bucket].Remove(e.node))
		return e.value
	}
	return nil
}

func (c *lfuCache) moveToNextBucket(e entry) {
	c.freqList[e.bucket].Remove(e.node)
	if e.bucket < c.capacity-1 {
		if e.bucket == c.minFreq && c.freqList[e.bucket].Len() == 0 {
			c.minFreq++
		}
		e.bucket++
	}
	c.freqList[e.bucket].PushFront(e.node)
}
