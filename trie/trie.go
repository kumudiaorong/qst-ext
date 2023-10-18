package trie

import (
	terr "qst-ext-appsearcher-go/error"
)

// isLeaf bool
// // 词频
// frequency int
// // 词
// word string
// }

// func NewTrieNode() *TrieNode {
// return &TrieNode{
// 	children:  make(map[rune]*TrieNode),
// 	isLeaf:    false,
// 	frequency: 0,
// 	word:      "",
// }
// }

// type Trie struct {
// root *TrieNode
// }

// func NewTrie() *Trie {
// return &Trie{
// 	root: NewTrieNode(),
// }
// }

// // 插入
// func (t *Trie) Insert(word string) {
// node := t.root
// for _, c := range word {
// 	if _, ok := node.children[c]; !ok {
// 		node.children[c] = NewTrieNode()
// 	}
// 	node = node.children[c]
// }
// node.isLeaf = true
// node.frequency++
// node.word = word
// }

// // 查找
// func (t *Trie) Search(word string) bool {
// node := t.root
// for _, c := range word {
// 	if _, ok := node.children[c]; !ok {
// 		return false
// 	}
// 	node = node.children[c]
// }
// return node.isLeaf
// }

// // 前缀查找
// func (t *Trie) StartsWith(prefix string) bool {
// node := t.root
// for _, c := range prefix {
// 	if _, ok := node.children[c]; !ok {
// 		return false
// 	}
// 	node = node.children[c]
// }
// return true
// }

// // 获取词频
// func (t *Trie) GetFrequency(word string) int {
// node := t.root
// for _, c := range word {
// 	if _, ok := node.children[c]; !ok {
// 		return 0
// 	}
// 	node = node.children[c]
// }
// return node.frequency
// }

// // 获取词
// func (t *Trie) GetWord(word string) string {
// node := t.root
// for _, c := range word {
// 	if _, ok := node.children[c]; !ok {
// 		return ""
// 	}
// 	node = node.children[c]
// }
// return node.word
// }

type TrieNode struct {
	// 子节点
	children map[rune]*TrieNode
	// 单词对应的值
	value interface{}
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
	}
}

const (
	// 未知错误
	SkipThis terr.Kind = iota + terr.MaxKind
)

func (tn *TrieNode) WalkChildren(fn func(*TrieNode) *terr.Error) (err *terr.Error) {
	for _, child := range tn.children {
		err = fn(child)
		if err != nil {
			switch err.Kind() {
			case terr.Unknown:
				return err
			case SkipThis:
				continue
			}
		}
		if err = child.WalkChildren(fn); err != nil {
			return
		}
	}
	return
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

// 插入
func (t *Trie) Insert(word string, value interface{}) {
	node := t.root
	for _, c := range word {
		child, ok := node.children[c]
		if !ok {
			child = NewTrieNode()
			node.children[c] = child
		}
		node = child
	}
	node.value = value
}

// 查找
func (t *Trie) Search(word string) interface{} {
	node := t.root
	for _, c := range word {
		if _, ok := node.children[c]; !ok {
			return nil
		}
		node = node.children[c]
	}
	return node.value
}

// 前缀查找
func (t *Trie) StartsWith(prefix string) (vs []interface{}) {
	vs = nil
	node := t.root
	for _, c := range prefix {
		if _, ok := node.children[c]; !ok {
			return
		}
		node = node.children[c]
	}
	vs = append(vs, node.value, node.WalkChildren(func(child *TrieNode) *terr.Error {
		vs = append(vs, child.value)
		return nil
	}))
	return
}

// 删除
func (t *Trie) Delete(word string) {
	node := t.root
	for _, c := range word {
		if _, ok := node.children[c]; !ok {
			return
		}
		node = node.children[c]
	}
	node.value = nil
}
