package timestamp

import (
	"time"

	"github.com/okneniz/parsec"
	. "github.com/okneniz/parsec/strings"
)

func dayOfWeekPrefix() parsec.Combinator[rune, Position, time.Weekday] {
	dwDict := map[string]time.Weekday{
		"Mon": time.Monday,
		"Tue": time.Tuesday,
		"Wed": time.Wednesday,
		"Thu": time.Thursday,
		"Fri": time.Friday,
		"Sat": time.Saturday,
		"Sun": time.Sunday,
	}

	return Padded(
		Try(Space("space")),
		MapStrings("expected day of week", dwDict),
	)
}

func monthPrefix() parsec.Combinator[rune, Position, time.Month] {
	monthDict := map[string]time.Month{
		"Jan": time.January,
		"Feb": time.February,
		"Mar": time.March,
		"Apr": time.April,
		"May": time.May,
		"Jun": time.June,
		"Jul": time.July,
		"Aug": time.August,
		"Sep": time.September,
		"Oct": time.October,
		"Nov": time.November,
		"Dec": time.December,
	}

	return Padded(
		Try(Space("space")),
		MapStrings("expected name of month", monthDict),
	)
}

func yearWithCentury() parsec.Combinator[rune, Position, int] {
	return Padded(
		Try(Space("space")),
		UnsignedN[int](4, "expected year with century"),
	)
}

func paddedDayNum() parsec.Combinator[rune, Position, int] {
	return Padded(
		Try(Space("space")),
		Choice(
			"expected day number",
			Try(UnsignedN[int](2, "expected day number")),
			UnsignedN[int](1, "expected day number"),
		),
	)
}

func paddedHourNum() parsec.Combinator[rune, Position, int] {
	return Padded(
		Try(Space("space")),
		UnsignedN[int](2, "expected hour number"),
	)
}

func paddedMinuteNum() parsec.Combinator[rune, Position, int] {
	return Padded(
		Try(Space("space")),
		UnsignedN[int](2, "expected minute number"),
	)
}

func paddedSecondNum() parsec.Combinator[rune, Position, int] {
	return Padded(
		Try(Space("space")),
		UnsignedN[int](2, "expected second number"),
	)
}
