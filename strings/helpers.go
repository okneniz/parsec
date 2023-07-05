package strings

import (
	p "git.sr.ht/~okneniz/parsec/common"
)

func Parens[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('('), body, Eq(')'))
}

func Braces[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('{'), body, Eq('}'))
}

func Angles[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('<'), body, Eq('>'))
}

func Squares[T any](
	body p.Combinator[rune, Position, T],
) p.Combinator[rune, Position, T] {
	return Between(Eq('['), body, Eq(']'))
}

func Semi() p.Combinator[rune, Position, rune] {
	return Eq(';')
}

func Comma() p.Combinator[rune, Position, rune] {
	return Eq(',')
}

func Colon() p.Combinator[rune, Position, rune] {
	return Eq(':')
}

func Dot() p.Combinator[rune, Position, rune] {
	return Eq('.')
}
