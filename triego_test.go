package triego

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

type word_test1 struct {
	w string // the word
	e bool   // if we expect the word in the trie or not
}

var word_tests1 = []word_test1{
	{"cat", true},
	{"doggy", true},
	{"dog", true},
	{"foo", false},
	{"can", false},
	{"c", false},
	{"cattelan", true},
}

func Test_trieFindsWords(t *testing.T) {
	rootTrie := NewTrie()

	for _, n := range word_tests1 {
		if n.e == true {
			rootTrie.AppendWord(n.w)
		}
		if has := rootTrie.HasWord(n.w); has != n.e {
			t.Errorf("Unexpected HasWord result for word '%s': got %v, expected %v", n.w, has, n.e)
			printTrie(rootTrie)
		}
	}
}

/*
 * A utility function to make sure
 * node append works properly for our trie
 */
func count_nodes(trie *Trie, i *int) {
	if len(trie.Children) == 0 {
		*i = *i + 1
		return
	}
	for _, v := range trie.Children {
		count_nodes(v, i)
	}

	*i = *i + 1
}

type node_test struct {
	words []string
	nodes int
}

var node_tests = []node_test{
	node_test{[]string{
		"romane",
		"romanus",
		"romulus",
		"rubens",
		"ruber",
		"rubicon",
		"rubicundus",
	}, 14,
	},
	node_test{[]string{
		"arma",
		"armatura",
		"armento",
	}, 5,
	},
}

func Test_trieNodeCount(t *testing.T) {
	for _, v := range node_tests {
		root_trie := NewTrie()
		// appending nodes
		for _, w := range v.words {
			root_trie.AppendWord(w)
		}

		var count int = 0
		count_nodes(root_trie, &count)
		if count != v.nodes {
			t.Errorf("Unexpected node count: got %d, expected %d", count, v.nodes)
			printTrie(root_trie)
		}
	}
}

type prefixes_test struct {
	words             []string
	query             string
	expected_prefixes []string
}

var prefixes_tests = []prefixes_test{
	{
		[]string{"dom", "domato", "domatore", "domenica", "domani"},
		"domat",
		[]string{"domato", "domatore"},
	},
	{
		[]string{"dopo domani", "domenica"},
		"domani",
		[]string{"dopo domani"},
	},
}

func Test_trieClosestWords(t *testing.T) {
	for _, v := range prefixes_tests {
		trie := NewTrie()
		for _, w := range v.words {
			trie.AppendWord(w)
		}

		prefixes := trie.ClosestWords(v.query) // []interface{} which are actually strings
		if len(prefixes) != len(v.expected_prefixes) {
			printTrie(trie)
			t.Errorf("Unexpected: expected prefixes length: %d, got: %d", len(v.expected_prefixes), len(prefixes))
		}

		for i := 0; i < len(v.expected_prefixes); i++ {
			found := false
			for j := 0; j < len(prefixes); j++ {
				if v.expected_prefixes[i] == prefixes[j].(string) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected: couldn't find expected prefix '%s'", v.expected_prefixes[i])
			}
		}
		if t.Failed() {
			for _, p := range prefixes {
				t.Log(string(p.([]rune)))
			}
			printTrie(trie)
		}
	}
}

type node_count_test struct {
	words          []string
	expected_count int
}

func Benchmark_nodesAllocation(b *testing.B) {
	b.StopTimer()

	// building country name
	// source from file
	file, err := os.Open("/usr/share/dict/words")
	words := make([]string, 0)
	if err != nil {
		b.Log("Cannot open expected file /usr/share/dict/words. Skipping this benchmark.")
		b.SkipNow()
		return
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		words = append(words, line[:len(line)-1])
	}

	b.Logf("Inserting %d words in the trie", b.N)

	rootTrie := NewTrie()

	b.ResetTimer()
	b.StartTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rootTrie.AppendWord(words[i%len(words)])
	}
}

func Benchmark_wordFind(b *testing.B) {
	b.StopTimer()

	// building country name
	// source from file
	file, err := os.Open("/usr/share/dict/words")
	words := make([]string, 0)
	if err != nil {
		b.Log("Cannot open expected file /usr/share/dict/words. Skipping this benchmark.")
		b.SkipNow()
		return
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		words = append(words, line[:len(line)-1])
	}
	rootTrie := NewTrie()

	for i := 0; i < b.N; i++ {
		rootTrie.AppendWord(words[i%len(words)])
	}

	b.ResetTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		found := rootTrie.HasWord(words[i%len(words)])
		if found != true {
			b.Fatalf("Unexpected: couldn't find word '%s'", words[i%len(words)])
		}
	}
}

/*
 * A few helper functions
 */
func printTrie(trie *Trie) {
	q := new_queue()
	last_depth := trie.depth

	q.enqueue(trie)
	for !q.is_empty() {
		n := q.dequeue()
		if n.isRoot {
			fmt.Print("/")
		} else {
			if n.depth > last_depth {
				fmt.Println()
			}
			fmt.Printf("(%s,%v,%d)\t", string(n.chars), n.IsWord, n.depth)
		}
		for _, c := range n.Children {
			q.enqueue(c)
		}
		last_depth = n.depth
	}
	fmt.Println()
}

type words_test struct {
	words []string
}

var words_tests = []words_test{
	{[]string{"a", "aaa", "abc", "zoro", "hephaestus"}},
}

func Test_Words(t *testing.T) {
	for _, v := range words_tests {
		trie := NewTrie()
		trie.AppendWords(v.words...)

		for _, w := range v.words {
			found := false
			for _, ws := range trie.Words() {
				if ws == w {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unable to find word '%s'", w)
			}
		}
		// debug
		if t.Failed() {
			printTrie(trie)
			t.Log(trie.Words())
		}
	}
}

type prefix_test_t struct {
	input_prefixes    []string
	expected_prefixes map[string]bool
}

var prefix_tests = []prefix_test_t{
	{
		[]string{"arma", "armatura", "armento"},
		map[string]bool{
			"arm":      false,
			"arma":     true,
			"armatura": true,
			"armento":  true,
		},
	},
}

func Test_EachPrefix(t *testing.T) {
	for _, tc := range prefix_tests {
		radix := NewTrie()
		for _, w := range tc.input_prefixes {
			radix.AppendWord(w)
		}

		var prefixes []string
		radix.EachPrefix(func(info PrefixInfo) (skip_subtree, halt bool) {
			v, ok := tc.expected_prefixes[info.Prefix]
			if ok == false || v != info.IsWord {
				t.Errorf(`Unexpected condition:
				For prefix '%s'
				* prefix expected: %v
				* prefix is word == (expected: %v, got: %v)
				`, info.Prefix, ok, v, info.IsWord)
				return false, false
			}

			if ok {
				prefixes = append(prefixes, info.Prefix)
			}
			return false, false
		})

		if len(prefixes) != len(tc.expected_prefixes) {
			t.Errorf("Unexpected count of prefixes: got %d, expected %d", len(prefixes), len(tc.expected_prefixes))
		}
	}
}
