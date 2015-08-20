package triego

func runes_eq(src, dst []rune) bool {
	if len(src) != len(dst) {
		return false
	}

	for i := range src {
		if src[i] != dst[i] {
			return false
		}
	}
	return true
}

/*
 * Returns the last index at which both
 * streams are still equal
 */
func same_until(src, dst []rune) (i int) {
	i = -1
	lsrc := len(src)
	ldst := len(dst)
	max  := lsrc

	if lsrc > ldst {
		max = ldst
	}

	for i = 0; i < max; i++ {
		if src[i] != dst[i] {
			break
		}
	}

	i--
	return
}

func max(args... int) (int) {
	if len(args) == 0 {
		panic("Cannot find max of empty list")
	}
	m := args[0]
	for _, v := range args {
		if v > m {
			m = v
		}
	}

	return m
}

func min(args... int) (int) {
	if len(args) == 0 {
		panic("Cannot find min of empty list")
	}
	m := args[0]
	for _, v := range args {
		if v < m {
			m = v
		}
	}

	return m
}

const q_PAGE_SIZE = 4096 // common page size

type queue struct {
	q []*Trie
	pages [][]*Trie
	h,t,page_index int
}

func (q *queue) enqueue(node *Trie) {
	if q.t == cap(q.q) {
		// moving to the next page
		q.page_index += 1

		// incrementing pages slice
		// if no empty pages are available
		if q.page_index == len(q.pages) {
			page := make([]*Trie, q_PAGE_SIZE)
			q.pages = append(q.pages, page)
		}
		q.q = q.pages[q.page_index]

		// resetting indexes
		q.t = 0
		q.h = 0
	}
	q.q[q.t] = node
	q.t += 1
}

func (q *queue) is_empty() (bool) {
	return q.h == q.t
}

func (q *queue) dequeue() (node *Trie) {
	if q.h == q.t {
		if q.page_index > 0 {
			q.page_index -= 1
			q.q = q.pages[q.page_index]
			q.h = 0
			q.t = len(q.q)
		} else {
			node     = nil
			return
		}
	}

	node = q.q[q.h]
	q.h += 1
	return
}

func (q *queue) clear() {
	q.h = 0
	q.t = 0
}

func new_queue() *queue {
	q := new(queue)
	q.q = make([]*Trie, q_PAGE_SIZE)
	q.pages = [][]*Trie{q.q}
	q.h = 0
	q.t = 0
	q.page_index = 0
	return q
}
