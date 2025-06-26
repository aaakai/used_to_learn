package test

import "fmt"

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

func ConstructorTrie() Trie {
	return Trie{root: NewTrieNode()}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = NewTrieNode()
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (t *Trie) Search(word string) bool {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			return false
		}
		node = node.children[ch]
	}
	return node.isEnd
}

func TrieTest() {
	trie := ConstructorTrie()
	trie.Insert("hello")
	trie.Insert("world")
	trie.Insert("hell")
	trie.Insert("help")
	trie.Insert("helping")
	trie.Insert("helpful")
	fmt.Println(trie.Search("hello"))   // 输出: true
	fmt.Println(trie.Search("world"))   // 输出: true
	fmt.Println(trie.Search("hell"))    // 输出: true
	fmt.Println(trie.Search("help"))    // 输出: true
	fmt.Println(trie.Search("helping")) // 输出: true
	fmt.Println(trie.Search("helpful")) // 输出: true
	fmt.Println(trie.Search("helps"))
	fmt.Println(trie.Search("hel"))
}
