package test

import (
	"container/list"
)

// Node 表示缓存中的每个节点
type Node struct {
	key   int
	value int
	freq  int
}

// LFUCache 结构体
type LFUCache struct {
	capacity int
	cache    map[int]*Node
	freqMap  map[int]*list.List
	minFreq  int
}

// Constructor 初始化 LFUCache
func NewConstructor(capacity int) LFUCache {
	return LFUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		freqMap:  make(map[int]*list.List),
		minFreq:  0,
	}
}

// Get 获取键的值
func (lfu *LFUCache) Get(key int) int {
	if node, exists := lfu.cache[key]; exists {
		lfu.update(node)
		return node.value
	}
	return -1
}

// Put 插入或更新键值对
func (lfu *LFUCache) Put(key int, value int) {
	if lfu.capacity <= 0 {
		return
	}

	if node, exists := lfu.cache[key]; exists {
		node.value = value
		lfu.update(node)
	} else {
		if len(lfu.cache) >= lfu.capacity {
			lfu.evict()
		}
		newNode := &Node{key: key, value: value, freq: 1}
		lfu.cache[key] = newNode
		if lfu.freqMap[1] == nil {
			lfu.freqMap[1] = list.New()
		}
		lfu.freqMap[1].PushFront(newNode)
		lfu.minFreq = 1
	}
}

// update 更新节点的频率
func (lfu *LFUCache) update(node *Node) {
	freq := node.freq
	lfu.freqMap[freq].Remove(getElement(lfu.freqMap[freq], node))

	if freq == lfu.minFreq && lfu.freqMap[freq].Len() == 0 {
		lfu.minFreq++
	}

	node.freq++
	if lfu.freqMap[node.freq] == nil {
		lfu.freqMap[node.freq] = list.New()
	}
	lfu.freqMap[node.freq].PushFront(node)
}

// evict 移除最不常使用的节点
func (lfu *LFUCache) evict() {
	if lfu.minFreq == 0 {
		return
	}

	list := lfu.freqMap[lfu.minFreq]
	oldest := list.Back()
	if oldest != nil {
		node := oldest.Value.(*Node)
		list.Remove(oldest)
		delete(lfu.cache, node.key)
		if list.Len() == 0 {
			delete(lfu.freqMap, lfu.minFreq)
		}
	}
}

// getElement 获取链表中的节点
func getElement(l *list.List, node *Node) *list.Element {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*Node) == node {
			return e
		}
	}
	return nil
}
