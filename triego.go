// This package exposes a simple Trie implementation
package triego

import (
	"github.com/alediaferia/stackgo"
	"unsafe"
	"strings"
)

const (
	k_DEFAULT_ALLOC_SIZE = 10
)

type Trie struct {
	IsWord   bool
	Parent   *Trie
	chars    []rune
	Children []*Trie
	isRoot   bool
	depth    int
}

type TrieNode Trie
type TriePtr *Trie

type TrieNodeIteratorCallback func(node *TrieNode, halt *bool)

type PrefixInfo struct {
	Prefix		  string
	IsWord		  bool
	Depth		  int

	// SharedLength gives
	// information about how
	// much length is shared with
	// the previous prefix,
	// aka, how many characters
	// of the previous prefix are shared
	// with the current one
	SharedLength  int
}
type PrefixIteratorCallback   func(PrefixInfo) (skip_subtree, halt bool)

// Initializes a new trie
func NewTrie() (t *Trie) {
	t = new(Trie)
	t.IsWord = false
	t.Parent = nil
	t.chars = make([]rune, 0, k_DEFAULT_ALLOC_SIZE)
	t.isRoot = true
	t.Children = make([]*Trie, 0)
	t.depth = 0

	return
}

// Returns true if this trie is root
func (t *Trie) IsRoot() bool {
	return t.isRoot
}

func (t *TrieNode) IsRoot() bool {
	return t.isRoot
}

// Returns the depth of the
// node within the whole trie
// it belongs to
func (t *TrieNode) Depth() int {
	return t.depth
}

// Appends a word to the trie
// This is a recursive function, so not that
// efficient.
func (t *Trie) AppendWord(word string) {
	t.append_radix([]rune(word))
}

func (t *Trie) AppendWords(words... string) {
	for _, w := range words {
		t.AppendWord(w)
	}
}

func (t *Trie) increase_depth() {
	last_depth := t.depth
	q 		   := new_queue()
	q.enqueue(t)

	for !q.is_empty() {
		n := q.dequeue()
		if n.depth > last_depth {
			last_depth = n.depth
		}

		n.depth++
		for _, c := range n.Children {
			q.enqueue(c)
		}
	}
}

func (t *Trie) delete_child(name string) {
	l := len(t.Children)
	for i := 0; i < l; i++ {
		if string(t.Children[i].chars) == name {
			if i == l - 1 {
				t.Children = t.Children[:i]
			} else {
				t.Children = append(t.Children[:i], t.Children[i + 1:]...)
			}
			return
		}
	}
}

// Appends the given suffix to
// the current radix.
// The algorithm follows the
// BFS traversal principles implemented
// iteratively. Given suffix is treated
// as a full word.
func (t *Trie) append_radix(suffix []rune) {
	cn        		 := t
	current_children := []*Trie{}
    var last_node *Trie = nil
    
    var e_range int = 0

	for cn != nil && len(suffix) > 0 {
 		if !cn.isRoot {
			if suffix[0] != cn.chars[0] {
				// the current subtree
				// is surely something we
				// don't want to examine further
				goto next
			}

			r := same_until(suffix, cn.chars)

			// case 1:
			// suffix and cn.chars completely
			// match: we found the node already
			// and we just have to make sure
			// the node is marked as word already
			if r == len(cn.chars) -1 && len(suffix) == len(cn.chars) {
				cn.IsWord = true
				return
			}

			if r > -1 {
				// storing this as the
				// last visited node
				last_node = cn

				// now adjusting our
				// given and temporary suffix
				// before proceeding:
				// they'll both start
				// at the next character
				suffix = suffix[r+1:]

				// finally storing the last
				// common index for both
				// suffixes
				e_range = r

				if r < len(cn.chars) - 1 {
					break
				}
			}
		}
		current_children = cn.Children
next:
		if len(current_children) != 0 {
			cn = current_children[0]
			current_children = current_children[1:]
		} else {
			cn = nil
		}
	}

	// no node found matching
	// part of the suffix we want
	// to append. A new one will
	// be created
	if last_node == nil {
		new_ := NewTrie()
		new_.isRoot = false
		new_.chars = make([]rune, len(suffix))
		copy(new_.chars, suffix)
		new_.Parent = t
		t.Children = append(t.Children, new_)
		new_.depth = t.depth + 1
		new_.IsWord = true
		return
	}

	// last_node now will contain the node
	// which constructs the closest match
	// to the suffix we are about to append
	//
	// we now need to split the matching node
	// content so that we can add our suffix

	// now adjusting current node
	// characters:
	// if current node is 'romane'
	// and we are about to append
	// word 'romanus' we want to preserve
	// just up to 'roman' and create 2 subnodes
	// 'e' and 'us' respectively
	orig_size := len(last_node.chars)

	left := last_node.chars[:e_range + 1] // will become the content of last_node
	sub1 := last_node.chars[e_range + 1:] // will become a new sub node
	sub2 := suffix // new sub node as well

	last_node.chars  = left

	was_word := last_node.IsWord

	if len(suffix) == 0 {
		last_node.IsWord = true
	} else if e_range +1 != orig_size {
		last_node.IsWord = false
	}

	// TODO: clarify this
	if len(sub1) != 0 {
		// appending sub1 contents
		sub1_c := new(Trie)
		sub1_c.isRoot = false
		sub1_c.IsWord = was_word
		sub1_c.Parent = last_node
		sub1_c.chars  = sub1
		sub1_c.depth  = last_node.depth // will increase this later
		sub1_c.Children = last_node.Children

		// we need to update children depth
		// since we have just moved this
		// subtree one lever lower
		sub1_c.increase_depth()

		// an important thing to remember is that
		// sub1_c inherits all the children from
		// last_node which has now been split
		last_node.Children = []*Trie{sub1_c}
	}


	if len(sub2) != 0 {
		// appending sub2 contents
		sub2_c := new(Trie)
		sub2_c.isRoot = false
		sub2_c.IsWord = true
		sub2_c.Parent = last_node
		sub2_c.chars  = sub2
		sub2_c.depth  = last_node.depth + 1
		sub2_c.Children = make([]*Trie, 0, 1)
		last_node.Children = append(last_node.Children, sub2_c)
	}
}

// Returns true if the word is found
// in the trie
func (t *Trie) HasWord(word string) bool {
	suffix           := []rune(word)
	cn        		 := t
	current_children := []*Trie{}

	for cn != nil {
		if !cn.isRoot {
			if suffix[0] != cn.chars[0] {
				goto next
			}
			last := same_until(suffix, cn.chars);
			if last == len(cn.chars) - 1 && len(suffix) == len(cn.chars) {
				return cn.IsWord // exact node match
			} else {
				suffix = suffix[last+1:]
			}
			// node contains the whole
			// suffix string: this means
			// we haven't found an exact
			// match
			if len(suffix) == 0 {
				return false
			}
		}

		current_children = cn.Children
next:
		if len(current_children) != 0 {
			cn = current_children[0]
			current_children = current_children[1:]
		} else {
			cn = nil
		}
	}

	return false
}

// Returns an array of closest words to
// the given word prefix string
func (t *Trie) ClosestWords(word string) []string {
	suffix           := []rune(word)
	cn        		 := t
	current_children := []*Trie{}
	var last_prefix_node *Trie = nil

	prefix := []rune{}

	for cn != nil {
		if !cn.isRoot {
			if suffix[0] != cn.chars[0] {
				goto next
			}
			last := same_until(suffix, cn.chars);
			// if the given suffix is equal
			// to the node we have an exact
			// match and therefore we return
			// just it if it is a word
			if last == len(cn.chars) - 1 && len(suffix) == len(cn.chars) {
				if cn.IsWord {
					return []string{string(cn.chars)}
				}
			}

			// if the given suffix
			// is still not an empty
			// string this means we have
			// found a prefix for the given
			// word
			if last > -1 {
				last_prefix_node = cn
				suffix = suffix[last+1:]
			}
			if len(suffix) == 0 {
				break
			} else {
				prefix = append(prefix, cn.chars[:last+1]...)
			}
		}

		current_children = cn.Children
next:
		if len(current_children) != 0 {
			cn = current_children[0]
			current_children = current_children[1:]
		} else {
			cn = nil
		}
	}

	if last_prefix_node != nil {
		suffixes := last_prefix_node.Words()
		for i := 0; i < len(suffixes); i++ {
			suffixes[i] = strings.Join([]string{string(prefix), suffixes[i]}, "")
		}
		return suffixes
	}

	return []string{}
}

// Returns a list with all the
// words present in the trie
func (t *Trie) Words() (words []string) {
	// DFS-based implementation for returning
	// all the words in the trie
	stack := stackgo.NewStack()

	words = make([]string, 0)
	word := make([]rune, 0)

	last_depth := 0
	last_append_size := 0

	stack.Push(unsafe.Pointer(t))
	for stack.Size() > 0 {
		node := TriePtr(stack.Pop().(unsafe.Pointer))

		if !node.isRoot {
			if node.depth <= last_depth {
				word = word[:len(word)-last_append_size]
			}

			word = append(word, node.chars...)

			if node.IsWord {
				words = append(words, string(word))
			}
		}

		for _, c := range node.Children {
			stack.Push(unsafe.Pointer(c))
		}
		last_depth = node.depth
		last_append_size = len(node.chars)
	}

	return
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

// Iterates for each prefix in the
// radix tree calling the given callback.
// The given callback can be used to
// guide the tree traversal.
func (t *Trie) EachPrefix(callback PrefixIteratorCallback) {
	stack  := stackgo.NewStack()
	prefix := []rune{}

	skipsubtree		  := false
	halt     		  := false
	added_lengths     := stackgo.NewStack()
	last_depth        := t.depth

	stack.Push(unsafe.Pointer(t))
	for stack.Size() != 0 {
		node := TriePtr(stack.Pop().(unsafe.Pointer))
		if !node.isRoot {
			// if we are now going up
			// in the radix (e.g. we have
			// finished with the current branch)
			// then we adjust the current prefix
			if last_depth >= node.depth {
				var length = 0
				for i := 0; i < (last_depth - node.depth) + 1; i++ {
					length += added_lengths.Pop().(int)
				}
				prefix = prefix[:len(prefix) - length]
			}
			last_depth = node.depth
			shared_length := len(prefix)
			prefix = append(prefix, node.chars...)
			added_lengths.Push(len(node.chars))

			// building the info
			// data to pass to the callback
			info := PrefixInfo{
				string(prefix),
				node.IsWord,
				node.depth,
				shared_length,
			}

			skipsubtree, halt = callback(info)
			if halt {
				return
			}
			if skipsubtree {
				continue
			}
		}

		for _, c := range node.Children {
			stack.Push(unsafe.Pointer(c))
		}
	}
}
