package triego
import "testing"

type runes_eq_test struct {
	src, dst	[]rune
	expectation bool
}

var runes_eq_tests = []runes_eq_test {
	{ []rune("cat"), []rune("cat"), true },
	{ []rune("foo"), []rune("bar"), false},
	{ []rune("foobar"), []rune("baz"), false },
	{ []rune("a"), []rune("a"), true},
}

func Test_runes_eq(t *testing.T) {
	for _, v := range runes_eq_tests {
		if eq := runes_eq(v.src, v.dst); eq != v.expectation {
			t.Errorf("runes_eq('%s', '%s'): got %v, expected: %v", string(v.src), string(v.dst), eq, v.expectation)
		}
	}
}

type same_until_test struct {
	src, dst 	[]rune
	expectation int
}

var same_until_tests = []same_until_test {
	{ []rune("foo"), []rune("foobar"), 2 },
	{ []rune("foo"), []rune("bar"), -1},
	{ []rune("bar"), []rune("baz"), 1},
	{ []rune("romulus"), []rune("romane"), 2},
}

func Test_same_until(t *testing.T) {
	for _, v := range same_until_tests {
		if eq := same_until(v.src, v.dst); eq != v.expectation {
			t.Errorf("same_until('%s', '%s'): got %v, expected: %v", string(v.src), string(v.dst), eq, v.expectation)
		}
	}
}

type queue_test struct {
	w string
}

var queue_tests = []queue_test {
	{"alessandro"},
	{"typeflow"},
	{"triego"},
}

func Test_queue(t *testing.T) {
	for _, v := range queue_tests {
		q := new_queue()
		word := make([]rune, 0, len(v.w))
		for _, c := range v.w {
			t_ := NewTrie()
			t_.chars = append(t_.chars, c)
			q.enqueue(t_)
		}

		for !q.is_empty() {
			t_ := q.dequeue()
			word = append(word, t_.chars[0])
		}

		if string(word) != v.w {
			t.Errorf("Unexpected characters dequeued: got '%s', expected '%s'", string(word), v.w)
		}
	}
}

func Benchmark_queue_enqueue(b *testing.B) {
	b.ReportAllocs()
	q := new_queue()
	t := NewTrie()
	b.StopTimer()
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		q.enqueue(t)
	}
}

func Benchmark_queue_dequeue(b *testing.B) {
	b.ReportAllocs()
	q := new_queue()
	t := NewTrie()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		q.enqueue(t)
	}
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		q.dequeue()
	}
}
