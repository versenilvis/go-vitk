// https://github.com/acomagu/trie
package tokenizer

import (
	"cmp"
	"fmt"
	"sort"
)

type Tree[K cmp.Ordered, V any] []node[K, V]

type node[K cmp.Ordered, V any] struct {
	value V
	next  int 
	label K
	match bool
	leaf  bool
}

type kv[K cmp.Ordered, V any] struct {
	k []K
	v V
}


func NewTrie[K cmp.Ordered, V any](keys [][]K, values []V) Tree[K, V] {
	if len(keys) != len(values) {
		panic("length mismatch of keys and values")
	}
	if len(keys) == 0 {
		return Tree[K, V]{node[K, V]{next: 1}}
	}

	kvs := make([]kv[K, V], 0, len(keys))
	for i, k := range keys {
		kvs = append(kvs, kv[K, V]{k, values[i]})
	}

	sort.Slice(kvs, func(i, j int) bool {
		a, b := kvs[i].k, kvs[j].k
		for i := 0; i < len(a) && i < len(b); i++ {
			if a[i] == b[i] {
				continue
			}
			return a[i] < b[i]
		}
		if len(a) == len(b) {
			panic(fmt.Sprintf("duplicate key detected: %v", kvs[i].k))
		}
		return len(a) < len(b)
	})

	t := Tree[K, V]{node[K, V]{next: 1}}
	t = t.construct(kvs, 0, 0)
	return t
}

func (t Tree[K, V]) construct(kvs []kv[K, V], depth, current int) Tree[K, V] {
	if depth == len(kvs[0].k) {
		t[current].match = true
		t[current].value = kvs[0].v
		kvs = kvs[1:]
		if len(kvs) == 0 {
			t[current].leaf = true
			return t
		}
	}

	p := []int{0}
	for i := 0; i < len(kvs); {
		t = append(t, node[K, V]{
			label: kvs[i].k[depth],
		})
		for c := kvs[i].k[depth]; i < len(kvs) && kvs[i].k[depth] == c; i++ {
		}
		p = append(p, i)
	}

	for i := 0; i < len(p)-1; i++ {
		t[t.nextOf(current)+i].next = len(t) - t.nextOf(current) - i
		t = t.construct(kvs[p[i]:p[i+1]], depth+1, t.nextOf(current)+i)
	}
	return t
}


func (t Tree[K, V]) Trace(path []K) Tree[K, V] {
	if len(t) == 0 {
		return nil
	}

	var u int
	for _, c := range path {
		if t[u].leaf {
			return nil
		}
		u = t.nextOf(u)
		v := t.nextOf(u)
		if v-u > 40 {
			u += sort.Search(v-u, func(m int) bool {
				return t[u+m].label >= c
			})
		} else {
			for ; u != v-1 && t[u].label < c; u++ {
			}
		}
		if u >= len(t) || t[u].label != c {
			return nil
		}
	}
	return t[u:]
}

func (t Tree[K, V]) TraceOne(c K) Tree[K, V] {
	if len(t) == 0 || t[0].leaf {
		return nil
	}
	u := t.nextOf(0)
	v := t.nextOf(u)

	if v-u > 40 {
		u += sort.Search(v-u, func(m int) bool {
			return t[u+m].label >= c
		})
	} else {
		for ; u != v-1 && t[u].label < c; u++ {
		}
	}
	if u >= len(t) || t[u].label != c {
		return nil
	}
	return t[u:]
}

func (t Tree[K, V]) Terminal() (V, bool) {
	var zero V
	if len(t) == 0 {
		return zero, false
	}
	return t[0].value, t[0].match
}

func (t Tree[K, V]) SearchLongestMatch(path []K, start int) (V, int) {
	var lastV V
	var lastLen int
	if start >= len(path) {
		return lastV, 0
	}

	curr := t.Trace(path[start : start+1])
	if curr == nil {
		return lastV, 0
	}

	if val, match := curr.Terminal(); match {
		lastV = val
		lastLen = 1
	}

	for i := start + 1; i < len(path); i++ {
		curr = curr.TraceOne(path[i])
		if curr == nil {
			break
		}
		if val, match := curr.Terminal(); match {
			lastV = val
			lastLen = i - start + 1
		}
	}
	return lastV, lastLen
}

func (t Tree[K, V]) nextOf(i int) int {
	return i + t[i].next
}
