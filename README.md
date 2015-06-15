# Triego

A simple trie implementation in Go.

# Usage

```go
package main

import (
  "github.com/alediaferia/triego"
  "fmt"
)

func countTrieNodes(trie *triegoTrie, i *int) {
	if len(trie.Children) == 0 {
		*i = *i + 1
		return
	}
	for _, v := range trie.Children {
		countTrieNodes(v, i)
	}

	*i = *i + 1
}

func main() {
    trie := triego.NewTrie()
    
    trie.AppendWord("trial")
    trie.AppendWord("trie")
    
    fmt.Println("Stored words:")
    for _, w := range trie.Words() {
        fmt.Println(w)
    }
    
    fmt.Println("")
    var nodes int = 0
    countTrieNodes(trie, &nodes)
    fmt.Printf("Number of allocated nodes: %d\n", nodes)
}

```

Output:
```
Stored words:
trie
trial

Number of allocated nodes: 7
```

# License
The code in this repository is released under the terms of the MIT license.
Copyright (c) Alessandro Diaferia <alediaferia@gmail.com>
