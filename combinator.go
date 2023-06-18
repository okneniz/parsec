package parsec


type Buffer[T any] interface {
	Read(bool) (T, bool)
	Seek(int)
	Position() int
}

type Combinator[T any, S any] func(Buffer[T]) (S, bool)

func Parse[T any, S any](buffer Buffer[T], parse Combinator[T, S]) (S, bool) {
	result, ok := parse(buffer)
	if ok {
		return result, ok
	}

	return *new(S), false
}

//

type bytesBuffer struct {
	data []byte
	position int
}

func (s *bytesBuffer) Read(x bool) (byte, bool) {
	if s.position >= len(s.data) {
		return 0, false
	}

	b := s.data[s.position]
	s.position++

	return b, true
}

func (s *bytesBuffer) Seek(x int) {
	s.position = x
}

func (s *bytesBuffer) Position() int {
	return s.position
}

func BytesBuffer(data []byte) *bytesBuffer {
	b := new(bytesBuffer)
	b.data = data
	b.position = 0
	return b
}

func ParseBytes[S any](data []byte, parse Combinator[byte, S]) (S, bool) {
	buf := BytesBuffer(data)
	return Parse[byte, S](buf, parse)
}

//

type Condition[T any] func(T) bool

type Composer[T any, S any, B any] func(T,S) B

func Satisfy[T any](greedy bool, f Condition[T]) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, bool) {
		token, ok := buffer.Read(greedy)
		if !ok {
			return *new(T), false
		}

		if f(token) {
			return token, true
		}

		return *new(T), false
	}
}

func Try[T any, S any](c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		pos := buffer.Position()

		r, ok := c(buffer)
		if !ok {
			buffer.Seek(pos)
		}

		return r, ok
	}
}

func Not[T any, S any](c Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		r, ok := c(buffer)
		if ok {
			return *new(S), false
		}

		return r, true
	}
}

func Or[T any, S any](x, y Combinator[T, S]) Combinator[T, S] {
	return func(buffer Buffer[T]) (S, bool) {
		r, ok := x(buffer)
		if !ok {
			return y(buffer)
		}

		return r, ok
	}
}

func And[T any, S any, B any, M any](
	x Combinator[T,S],
	y Combinator[T,B],
	compose Composer[S,B,M],
) Combinator[T,M] {
	return func(buffer Buffer[T]) (M, bool) {
		first, ok := x(buffer)
		if !ok {
			return *new(M), false
		}

		second, ok := y(buffer)
		if !ok {
			return *new(M), false
		}

		return compose(first, second), true
	}
}

func Token[T comparable](greedy bool, t T) Combinator[T, T] {
	return func(buffer Buffer[T]) (T, bool) {
		r, ok := buffer.Read(greedy)
		if !ok || (r != t) {
			return *new(T), false
		}

		return r, true
	}
}

func Slice[T comparable, S any](cs ...Combinator[T,S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, len(cs))

		for _, c := range cs {
			t, ok := c(buffer)
			if !ok {
				return nil, false
			}

			result = append(result, t)
		}

		return result, true
	}
}

func Optional[T any, S any](greedy bool, c Combinator[T, S]) Combinator[T, *S] {
	return func(buffer Buffer[T]) (*S, bool) {
		t, ok := c(buffer)
		if ok {
			return &t, ok
		}

		return nil, true
	}
}

func Many[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		result := make([]S, 0, cap)

		for {
			t, ok := c(buffer)
			if !ok {
				break
			}

			result = append(result, t)
		}

		return result, true
	}
}

func Some[T any, S any](cap int, c Combinator[T, S]) Combinator[T, []S] {
	return func(buffer Buffer[T]) ([]S, bool) {
		cc := Many(cap, c)

		t, ok := cc(buffer)
		if !ok {
			return t, ok
		}
		if len(t) > 0 {
			return t, true
		}

		return nil, false
	}
}

func Before[T any, S any, B any, Z any](
	pre Combinator[T,S],
	body Combinator[T,B],
	compose Composer[S,B,Z],
) Combinator[T, Z] {
	return func(buffer Buffer[T]) (Z, bool) {
		prefix, ok := pre(buffer)
		if !ok {
			return *new(Z), false
		}

		suffix, ok := body(buffer)
		if !ok {
			return *new(Z), false
		}

		return compose(prefix, suffix), true
	}
}

func flip[T,S,B any](c func(T,S) B) func(S,T)B {
	return func(s S, t T) B {
		return c(t, s)
	}
}

func After[T any, S any, B any, Z any](
	pre Combinator[T,S],
	body Combinator[T,B],
	compose Composer[S,B,Z],
) Combinator[T, Z] {
	return Before(body, pre, flip(compose))
}

type Composer3[T, S, B, M any] func(T, S, B) M

func Between[T any, S any, B any, M any, Z any](
	pre Combinator[T,S],
	c Combinator[T,B],
	suf Combinator[T,M],
	compose Composer3[S,B,M,Z],
) Combinator[T, Z] {
	return func(buffer Buffer[T]) (Z, bool) {
		prefix, ok := pre(buffer)
		if !ok {
			return *new(Z), false
		}

		body, ok := c(buffer)
		if !ok {
			return *new(Z), false
		}

		suffix, ok := suf(buffer)
		if !ok {
			return *new(Z), false
		}

		return compose(prefix, body, suffix), true
	}
}
