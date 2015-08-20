# Triego

Despite the name the internal of this data structure is implemented
using a [Radix Tree](https://en.wikipedia.org/wiki/Radix_tree). 

This project is one of the funding blocks of [typeflow-webservice](https://github.com/typeflow/typeflow-webservice) but,
although mostly designed around it, it features a rather general purpose API around
a Radix Tree.

Implementation considerations
-----------------------------

The whole implementation has been designed with prefix iteration efficiency in mind.
This implies that feeding the radix tree is way slower than a simple trie, but looking for words
or prefixes is way faster due to the reduced amount of nodes to traverse.
So, although insertion should be still `MAX(O(N*Log(N)))`, the algorithm requires some additional time mostly
due to memory allocations and/or resizing of existing nodes.

Basic usage
-----------

```go
package main

import (
  "github.com/alediaferia/triego"
  "fmt"
)

func countTrieNodes(trie *triego.Trie, i *int) {
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

Number of allocated nodes: 4
```

That is, 3 actual nodes plus the root node. This is because 1 node is required for the "**tri**" prefix
and just 2 additional nodes for "**e**" and "**al**":

* tri
  - e
  - al


Advanced concepts
-----------------

### Prefixes iteration

The provided API can be used to iterate over all the prefixes in the radix.

You need to provide the `EachPrefix` API with a callback having this signature:

```go
type PrefixIteratorCallback   func(PrefixInfo) (skip_subtree, halt bool)
```

```PrefixInfo``` holds information about the current prefix. You can find out [here](/triego.go#L27-L39).
The callback can also be used to skip entire subtrees making the search even faster if it meets
certain conditions of your choice.
To do so, simply return ```true``` as the first return argument.

```go

radix.EachPrefix(func(info PrefixInfo) (skip_subtree, halt bool) {
    value := some_computation(info.Prefix)
    
    if value > MAX_THRESHOLD {
        // this value is too high
        // and we are not interested in
        // further prefixes containing
        // the current one so we skip
        // the current subtree altogether
        return true, false
    }
    
    ...
    
    return false, false
})

```

# License
The code in this repository is released under the terms of the MIT license.
Copyright (c) Alessandro Diaferia <alediaferia@gmail.com>
