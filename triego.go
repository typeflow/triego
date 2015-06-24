// This package exposes a simple Trie implementation
package triego

import (
	"github.com/alediaferia/stackgo"
	"unsafe"
)

type Trie struct {
	IsWord   bool
	Parent   *Trie
	C        rune
	Children map[rune]*Trie
	isRoot   bool
}

type TrieNode Trie
type TriePtr *Trie

type TrieNodeIteratorCallback func(node *TrieNode, halt *bool) ()

// Initializes a new trie
func NewTrie() (t *Trie) {
	t = new(Trie)
	t.IsWord = false
	t.Parent = nil
	t.C = 0
	t.isRoot = true
	t.Children = make(map[rune](*Trie))

	return
}

// Returns true if this trie is root
func (t *Trie) IsRoot() bool {
	return t.isRoot
}

// Appends a word to the trie
// This is a recursive function, so not that
// efficient.
func (t *Trie) AppendWord(word string) {
	t.append([]rune(word), true)
}

func (t *Trie) append(suffix []rune, makesWord bool) {
	if len(suffix) == 0 {
		return
	}

	// if there is already a node
	// holding this character we
	// move forward and append
	// the remaining part
	if c,ok := t.Children[suffix[0]]; ok {
		c.append(suffix[1:], makesWord)
		return
	}

	tc := NewTrie()
	tc.Parent = t
	t.Children[suffix[0]] = tc
	tc.C = suffix[0]
	tc.isRoot = false

	if len(suffix) > 1 {
		tc.append(suffix[1:], makesWord)
	} else {
		tc.IsWord = makesWord
	}
}


// Returns true if the word is present
// in the trie
func (t *Trie) HasWord(word string) bool {
	currentSlice := []rune(word)
	currentRoot  := t

	for len(currentSlice) > 0 {
		c, ok := currentRoot.Children[currentSlice[0]]
		if len(currentSlice) == 1 && ok == true && c.IsWord {
			return true
		} else if !ok {
			return false
		}
		currentSlice = currentSlice[1:]
		currentRoot  = c
	}

	return false
}

// Returns a list with all the
// words present in the trie
func (t *Trie) Words() (words []string) {
    // DFS-based implementation for returning
	// all the words in the trie
	stack := stackgo.NewStack()
	node := t

	words = make([]string, 0)
	word  := make([]rune, 0)

	stack.Push(unsafe.Pointer(node))
	for stack.Size() > 0 {
		node = TriePtr(stack.Pop().(unsafe.Pointer))
		if !node.isRoot {
			word = append(word, node.C)
		}

		if len(node.Children) == 0 {
			if node.IsWord {
				words = append(words, string(word))
			}
			word = word[:len(word) - 1]
		}
		for _, c := range node.Children {
			stack.Push(unsafe.Pointer(c))
		}
	}

	return
}

func (t *TrieNode) Character() rune {
	return t.C
}

func (t *Trie) EachNode(callback TrieNodeIteratorCallback) {
	// still a DFS-based implementation
	stack := stackgo.NewStack()
	node := t

	stack.Push(unsafe.Pointer(node))

	stop := false
	for stack.Size() > 0 {
		node = TriePtr(stack.Pop().(unsafe.Pointer))

		callback((*TrieNode)(node), &stop)

		if stop == true {
			return
		}

		for _, c := range node.Children {
			stack.Push(unsafe.Pointer(c))
		}
	}
}
