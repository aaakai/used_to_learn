package test

import "container/list"

type ValueNode struct {
	key   int
	value int
	freq  int
	listE *list.Element
}

type LFUCacheNew struct {
	capacity int
	minFreq  int
	cache    map[int]*ValueNode
	freqMap  map[int]*list.List
}

func NewLFUCacheNew(capacity int) *LFUCacheNew {
	return &LFUCacheNew{
		capacity: capacity,
		minFreq:  0,
		cache:    make(map[int]*ValueNode),
		freqMap:  make(map[int]*list.List),
	}
}

func (lfu *LFUCacheNew) Get(key int) int {
	if node, ok := lfu.cache[key]; ok {
		lfu.update(node)
		return node.value
	}
	return -1
}

func (lfu *LFUCacheNew) Put(key int, value int) {
	if lfu.capacity <= 0 {
		return
	}
	if node, ok := lfu.cache[key]; ok {
		node.value = value
		lfu.update(node)
		return
	} else {
		if len(lfu.cache) >= lfu.capacity {
			lfu.removeMinFreq()
		}
		newNode := &ValueNode{
			key:   key,
			value: value,
			freq:  1,
		}
		lfu.cache[key] = newNode
		if lfu.freqMap[1] == nil {
			lfu.freqMap[1] = list.New()
		}
		newNode.listE = lfu.freqMap[1].PushFront(newNode)
		lfu.minFreq = 1
	}
}

func (lfu *LFUCacheNew) update(node *ValueNode) {
	fre := node.freq
	lfu.freqMap[fre].Remove(node.listE)
	if fre == lfu.minFreq && lfu.freqMap[fre].Len() == 0 {
		lfu.minFreq++
	}
	node.freq++
	if lfu.freqMap[node.freq] == nil {
		lfu.freqMap[node.freq] = list.New()
	}
	node.listE = lfu.freqMap[node.freq].PushFront(node)
}

func (lfu *LFUCacheNew) removeMinFreq() {
	if lfu.minFreq == 0 {
		return
	}
	minList := lfu.freqMap[lfu.minFreq]
	remove := minList.Back()
	if remove != nil {
		delete(lfu.cache, remove.Value.(*ValueNode).key)
		minList.Remove(remove)
		if minList.Len() == 0 {
			delete(lfu.freqMap, lfu.minFreq)
		}
	}
}
