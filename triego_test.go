package triego

import (
	"testing"
	/*"fmt"*/
)

func Test_trieFindsWords(t *testing.T) {
	rootTrie := NewTrie()

	// manually appending two short
	// words to make sure test is
	// exclusively run against
	// the find function
	//
	// words added: cat, dog

	// dog
	tr := NewTrie()
	tr.isRoot = false
	tr.Parent = rootTrie
	tr.C = 'd'

	rootTrie.Children['d'] = tr

	tr1 := NewTrie()
	tr1.isRoot = false
	tr1.Parent = tr
	tr1.C = 'o'
	tr.Children['o'] = tr1

	tr2 := NewTrie()
	tr2.isRoot = false
	tr2.Parent = tr1
	tr2.C = 'g'
	tr2.IsWord = true
	tr1.Children['g'] = tr2

	// cat
	tr3 := NewTrie()
	tr3.isRoot = false
	tr3.Parent = rootTrie
	tr3.C = 'c'

	rootTrie.Children['c'] = tr3

	tr4 := NewTrie()
	tr4.isRoot = false
	tr4.Parent = tr3
	tr4.C = 'a'
	tr3.Children['a'] = tr4

	tr5 := NewTrie()
	tr5.isRoot = false
	tr5.Parent = tr4
	tr5.C = 't'
	tr5.IsWord = true
	tr4.Children['t'] = tr5

	if rootTrie.HasWord([]rune("dog")) == false {
		t.Errorf("Finding word 'dog' in trie fails")
	}

	if rootTrie.HasWord([]rune("cat")) == false {
		t.Errorf("Finding word 'cat' in trie fails")
	}

	if rootTrie.HasWord([]rune("foo")) == true {
		t.Errorf("Finding word 'foo' in trie unexpectedly succeeds")
	}

	var i int = 0
	countTries(rootTrie, &i)
	if i != 7 {
		t.Fatalf("Expected 7 nodes, got %d", i)
	}
}

/*
 * A utility function to make sure
 * node append workd properly for our trie
 */
func countTries(trie *Trie, i *int) {
	if len(trie.Children) == 0 {
		*i = *i + 1
		return
	}
	for _, v := range trie.Children {
		countTries(v, i)
	}

	*i = *i + 1
}

func Test_trieAppendsWords(t *testing.T) {
	rootTrie := NewTrie()

	const (
		w1 = "testWord1"
		w2 = "testWord2"
	)

	expectedWords := [...]string{w1, w2}

	rootTrie.AppendWord(w1)
	rootTrie.AppendWord(w2)

	if rootTrie.HasWord([]rune(w1)) == false {
		t.Errorf("Finding word '%s' in trie fails", w1)
	}

	if rootTrie.HasWord([]rune(w2)) == false {
		t.Errorf("Finding word '%s' in trie fails", w2)
	}
	var i int = 0
	countTries(rootTrie, &i)
	if i != 11 {
		t.Fatalf("Expected 11 nodes, got %d", i)
	}

	words := rootTrie.Words()
	for _, word := range words {
		found := false
		for _, expectedWord := range expectedWords {
			if word == expectedWord {
				found = true
			}
		}

		if !found {
			t.Fatalf("Cannot find expected words in the list of words in the trie: '%s' not in %v", word, expectedWords)
		}
	}
}

/*
 * A few helper functions
 */
/*func printTrie(trie *Trie) {
	for _, v := range trie.Children {
		runes := make([]rune, 0)
		printTrie_(v, append(runes, v.C))
	}
}

func printTrie_(trie *Trie, runes []rune) {
	for _,v := range trie.Children {
		if v.IsWord {
			fmt.Println(string(append(runes, v.C)))
		}
		printTrie_(v, append(runes, v.C))
	}
}*/
