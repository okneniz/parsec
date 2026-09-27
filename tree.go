package parsec

// Tree is a lookup structure which matches the current buffer content
// against a set of keys and returns the combinator stored for
// the matched key.
type Tree[T any, P any, O any, S any] interface {
	Lookup(state S, buf Buffer[T, P]) (Combinator[T, P, O, S], Error[P])
}

type tree[T comparable, P any, K comparable, O any, S any] struct {
	children map[T]*tree[T, P, K, O, S]
	value    Combinator[T, P, O, S]
}

var _ Tree[rune, int, float32, Stateless] = new(tree[rune, int, string, float32, Stateless])

// NewLongestPrefixTree builds a trie from the cases keys, each key split
// into items with split, with the corresponding combinators stored
// in the terminal nodes. Lookup consumes input while it matches some key
// and returns the combinator of the longest matched prefix,
// or nil when the first item already matches nothing.
func NewLongestPrefixTree[T comparable, P any, K comparable, O any, S any](
	cases map[K]Combinator[T, P, O, S],
	split func(K) []T,
) Tree[T, P, O, S] {
	root := new(tree[T, P, K, O, S])
	root.children = make(map[T]*tree[T, P, K, O, S])

	for key, value := range cases {
		seq := split(key)
		current := root

		for _, x := range seq {
			child, exists := current.children[x]
			if !exists {
				child = &tree[T, P, K, O, S]{
					children: make(map[T]*tree[T, P, K, O, S]),
				}

				current.children[x] = child
			}

			current = child
		}

		current.value = value
	}

	return root
}

func (tree *tree[T, P, K, O, S]) Lookup(state S, buf Buffer[T, P]) (Combinator[T, P, O, S], Error[P]) {
	current := tree.children

	var longestPrefix Combinator[T, P, O, S]

	for len(current) > 0 {
		pos := buf.Position()

		x, err := buf.Read(true)
		if err != nil {
			if seekErr := buf.Seek(pos); seekErr != nil {
				return nil, NewParseError(pos, seekErr.Error())
			}

			break
		}

		next, exists := current[x]
		if !exists {
			if seekErr := buf.Seek(pos); seekErr != nil {
				return nil, NewParseError(pos, seekErr.Error())
			}

			break
		}

		if next.value != nil {
			longestPrefix = next.value
		}

		current = next.children
	}

	return longestPrefix, nil
}
