package trie

import (
	"testing"
)

func TestTrie(t *testing.T) {
	tr := NewTrie()
	tr.Insert("hello", "hello")
	tr.Insert("world", 1)
	tr.Insert("好", 2)
	tr.Insert("好不", 3)

	i := tr.Search("hello")
	if i != "hello" {
		t.Error("search hello error")
	}
	i = tr.Search("world")
	if i != 1 {
		t.Error("search world error")
	}
	i = tr.Search("好")
	if i != 2 {
		t.Error("search 好 error")
	}
	a := tr.StartsWith("好")
	for _, v := range a {
		t.Log(v)
	}

}
