package timestamp

import (
	"time"

	"github.com/okneniz/parsec"
	. "github.com/okneniz/parsec/strings"
)

func dayOfWeekPrefix() parsec.Combinator[rune, Position, time.Weekday, parsec.Stateless] {
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
		Try(Space[parsec.Stateless]("space")),
		MapStrings[time.Weekday, parsec.Stateless]("expected day of week", dwDict),
	)
}

func monthPrefix() parsec.Combinator[rune, Position, time.Month, parsec.Stateless] {
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
		Try(Space[parsec.Stateless]("space")),
		MapStrings[time.Month, parsec.Stateless]("expected name of month", monthDict),
	)
}

func yearWithCentury() parsec.Combinator[rune, Position, int, parsec.Stateless] {
	return Padded(
		Try(Space[parsec.Stateless]("space")),
		UnsignedN[int, parsec.Stateless](4, "expected year with century"),
	)
}

func paddedDayNum() parsec.Combinator[rune, Position, int, parsec.Stateless] {
	return Padded(
		Try(Space[parsec.Stateless]("space")),
		Choice(
			"expected day number",
			Try(UnsignedN[int, parsec.Stateless](2, "expected day number")),
			UnsignedN[int, parsec.Stateless](1, "expected day number"),
		),
	)
}

func paddedHourNum() parsec.Combinator[rune, Position, int, parsec.Stateless] {
	return Padded(
		Try(Space[parsec.Stateless]("space")),
		UnsignedN[int, parsec.Stateless](2, "expected hour number"),
	)
}

func paddedMinuteNum() parsec.Combinator[rune, Position, int, parsec.Stateless] {
	return Padded(
		Try(Space[parsec.Stateless]("space")),
		UnsignedN[int, parsec.Stateless](2, "expected minute number"),
	)
}

func paddedSecondNum() parsec.Combinator[rune, Position, int, parsec.Stateless] {
	return Padded(
		Try(Space[parsec.Stateless]("space")),
		UnsignedN[int, parsec.Stateless](2, "expected second number"),
	)
}
