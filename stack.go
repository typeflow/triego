package triego

// This is mostly a copy of  http://github.com/alediaferia/stackgo
// with the static type *Trie for performance reasons
// for taking advantage of the varargs-based push
//
// For the original visit http://github.com/alediaferia/stackgo

type Stack struct {
	size int
	currentPage []*Trie
	pages [][]*Trie
	offset int
	capacity int
	pageSize int
	currentPageIndex int
}

const s_DefaultAllocPageSize = 4096

func NewStack() *Stack {
	stack := new(Stack)
	stack.currentPage = make([]*Trie, s_DefaultAllocPageSize)
	stack.pages = [][]*Trie{stack.currentPage}
	stack.offset = 0
	stack.capacity = s_DefaultAllocPageSize
	stack.pageSize = s_DefaultAllocPageSize
	stack.size = 0
	stack.currentPageIndex = 0

	return stack
}

func (s *Stack) Push(elem... *Trie) {
	if elem == nil || len(elem) == 0 {
		return
	}

	if s.size == s.capacity {
		pages_count := len(elem) / s.pageSize
		if len(elem) % s.pageSize != 0 {
			pages_count++
		}
		s.capacity += s.pageSize

		s.currentPage = make([]*Trie, s.pageSize)
		s.pages = append(s.pages, s.currentPage)
		s.currentPageIndex++

		pages_count--
		for pages_count > 0 {
			page := make([]*Trie, s.pageSize)
			s.pages = append(s.pages, page)
		}

		s.offset = 0
	}

	available := len(s.currentPage) - s.offset
	for len(elem) > available {
		copy(s.currentPage[s.offset:], elem[:available])
		s.currentPage = s.pages[s.currentPageIndex + 1]
		s.currentPageIndex++
		elem = elem[available:]
		s.offset = 0
		available = len(s.currentPage)
	}

	copy(s.currentPage[s.offset:], elem)
	s.offset += len(elem)
	s.size += len(elem)
}

func (s *Stack) Pop() (elem *Trie) {
	if s.size == 0 {
		return nil
	}

	s.offset--
	s.size--
	if s.offset < 0 {
		s.offset = s.pageSize - 1

		s.currentPage, s.pages = s.pages[len(s.pages) - 2], s.pages[:len(s.pages) - 1]
		s.capacity -= s.pageSize
		s.currentPageIndex--
	}

	elem = s.currentPage[s.offset]

	return
}

func (s *Stack) Top() (elem *Trie) {
	if s.size == 0 {
		return nil
	}

	off := s.offset - 1
	if off < 0 {
		page := s.pages[len(s.pages)-1]
		elem = page[len(page)-1]
		return
	}
	elem = s.currentPage[off]
	return
}

func (s *Stack) Size() int {
	return s.size
}

