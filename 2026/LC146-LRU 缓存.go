package _026

type LRUCache struct {
	size       int
	cap        int
	cache      map[int]*DLinkedNode
	head, tail *DLinkedNode
}

type DLinkedNode struct {
	key, value int
	prev, next *DLinkedNode
}

func initDLinkedNode(key, value int) *DLinkedNode {
	return &DLinkedNode{
		prev:  nil,
		next:  nil,
		key:   key,
		value: value,
	}
}

func Constructor(capacity int) LRUCache {
	lRUCache := LRUCache{
		cap:   capacity,
		size:  0,
		cache: make(map[int]*DLinkedNode),
		head:  initDLinkedNode(0, 0),
		tail:  initDLinkedNode(0, 0),
	}
	lRUCache.head.next = lRUCache.tail
	lRUCache.tail.prev = lRUCache.head
	return lRUCache
}

func (this *LRUCache) Get(key int) int {
	if _, ok := this.cache[key]; !ok {
		return -1
	}
	node := this.cache[key]
	this.moveToHead(node)
	return node.value
}

func (this *LRUCache) Put(key int, value int) {
	if _, ok := this.cache[key]; !ok {
		node := initDLinkedNode(key, value)
		this.addToHead(node)
		this.size++
		if this.size > this.cap {
			delNode := this.removeTail()
			delete(this.cache, delNode.key)
			this.size--
		}
	} else {
		node := this.cache[key]
		node.value = value
		this.moveToHead(node)
	}
}

func (this *LRUCache) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (this *LRUCache) addToHead(node *DLinkedNode) {
	node.prev = this.head
	node.next = this.head.next
	this.head.next.prev = node
	this.head.next = node
}

func (this *LRUCache) moveToHead(node *DLinkedNode) {
	this.removeNode(node)
	this.addToHead(node)
}

func (this *LRUCache) removeTail() *DLinkedNode {
	node := this.tail.prev
	this.removeNode(node)
	return node
}
