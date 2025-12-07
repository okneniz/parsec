package common

type Tree[T any, P any, S any] interface {
	Lookup(Buffer[T, P]) (Combinator[T, P, S], Error[P])
}

type tree[T comparable, P any, K comparable, V any] struct {
	children map[K]*tree[T, P, K, V]
	value    Combinator[K, P, V]
}

var _ Tree[rune, int, float32] = new(tree[string, int, rune, float32])

func NewLongestPrefixTree[T comparable, P any, K comparable, V any](
	cases map[T]Combinator[K, P, V],
	split func(T) []K,
) Tree[K, P, V] {
	root := new(tree[T, P, K, V])
	root.children = make(map[K]*tree[T, P, K, V])

	for key, value := range cases {
		seq := split(key)
		current := root

		for _, x := range seq {
			// TODO : check and handle conflicts?

			child, exists := current.children[x]
			if !exists {
				child = &tree[T, P, K, V]{
					children: make(map[K]*tree[T, P, K, V]),
				}

				current.children[x] = child
			}

			current = child
		}

		current.value = value
	}

	return root
}

func (tree *tree[T, P, K, V]) Lookup(buf Buffer[K, P]) (Combinator[K, P, V], Error[P]) {
	current := tree.children

	var longestPrefix Combinator[K, P, V]

	for len(current) > 0 {
		pos := buf.Position()

		x, err := buf.Read(true)
		if err != nil {
			if seekErr := buf.Seek(pos); seekErr != nil {
				return nil, NewParseError(pos, err.Error())
			}

			break
		}

		next, exists := current[x]
		if !exists {
			if seekErr := buf.Seek(pos); seekErr != nil {
				return nil, NewParseError(pos, err.Error())
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
