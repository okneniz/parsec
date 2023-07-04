package bytes

import (
	"git.sr.ht/~okneniz/parsec/common"
)

type Combinator[T any] common.Combinator[byte, int, T]

type Condition common.Condition[byte]
