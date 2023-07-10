package strings

type trie[K comparable, V any] struct {
	children map[K]*node[K, V]
}

type node[K comparable, V any] struct {
	end      bool
	value    V
	children map[K]*node[K, V]
}

func stringTrie[V any](cases map[string]V) *trie[rune, V] {
	t := new(trie[rune, V])
	t.children = make(map[rune]*node[rune, V])

	for cs, value := range cases {
		current := t.children

		for i, r := range cs {
			next, exists := current[r]
			if !exists {
				next = &node[rune, V]{
					children: make(map[rune]*node[rune, V]),
					end:      i == len(cs)-1,
					value:    value,
				}

				current[r] = next
			}

			current = next.children
		}
	}

	return t
}
