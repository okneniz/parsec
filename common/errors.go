package common

type Error string

const (
	EndOfFile         = Error("end of file")
	NotEnoughElements = Error("not enough of elements")
	NothingMatched    = Error("nothing matched")
)

func (s Error) Error() string {
	return string(s)
}
