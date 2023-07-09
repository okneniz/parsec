package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Eq(t rune) p.Combinator[rune, Position, rune] {
	return p.Eq[rune, Position](t)
}

func NotEq(t rune) p.Combinator[rune, Position, rune] {
	return p.NotEq[rune, Position](t)
}

func OneOf(data ...rune) p.Combinator[rune, Position, rune] {
	return p.OneOf[rune, Position](data...)
}

func NoneOf(data ...rune) p.Combinator[rune, Position, rune] {
	return p.NoneOf[rune, Position](data...)
}

func SequenceOf(data ...rune) p.Combinator[rune, Position, []rune] {
	return p.SequenceOf[rune, Position](data...)
}

type trie[K comparable, V any] struct {
	children map[K]*node[K, V]
}

type node[K comparable, V any] struct {
	end bool
	value V
	children map[K]*node[K,V]
}

func stringTrie[V any](cases map[string]V) *trie[rune,V] {
	t := new(trie[rune, V])
	t.children = make(map[rune]*node[rune, V])

	for cs, value := range cases {
		current := t.children

		for i, r := range cs {
			next, exists := current[r]
			if !exists {
				next = &node[rune,V]{
					children: make(map[rune]*node[rune,V]),
					end: i == len(cs) - 1,
					value: value,
				}

				current[r] = next
			}

			current = next.children
		}
	}

	return t
}

func MapStrings[V any](
	cases map[string]V,
) p.Combinator[rune, Position, V] {
	tr := stringTrie(cases)

	return func(buffer p.Buffer[rune, Position]) (V, error) {
		current := tr.children
		pos := buffer.Position()

		var result *V

		for {
			r, err := buffer.Read(true)
			if err != nil {
				break
			}

			next, exists := current[r]
			if !exists {
				break
			}

			if next.end {
				result = &next.value
				pos = buffer.Position()
			}

			current = next.children
		}

		buffer.Seek(pos)

		if result == nil {
			return *new(V), p.NothingMatched
		}

		return *result, nil
	}
}


func Map[K comparable, V any](
	cases map[K]V,
	c p.Combinator[rune, Position, K],
) p.Combinator[rune, Position, V] {
	return p.Map[rune, Position, K, V](cases, c)
}

func String(str string) p.Combinator[rune, Position, string] {
	return func(buffer p.Buffer[rune, Position]) (string, error) {
		for _, r := range str {
			c, err := buffer.Read(true)
			if err != nil {
				return "", err
			}

			if r != c {
				return "", p.NothingMatched
			}
		}

		return str, nil
	}
}

func OneOfStrings(strs ...string) p.Combinator[rune, Position, string] {
	combs := make([]p.Combinator[rune, Position, string], len(strs))

	for i, str := range strs {
		combs[i] = Try(String(str))
	}

	return Choice(combs...)
}
