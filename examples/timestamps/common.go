package timestamp

import (
	"time"

	p "git.sr.ht/~okneniz/parsec/common"
	. "git.sr.ht/~okneniz/parsec/strings"
)

func dayOfWeekPrefix() p.Combinator[rune, Position, time.Weekday] {
	dwDict := map[string]time.Weekday{
		"Mon": time.Monday,
		"Tue": time.Tuesday,
		"Wed": time.Wednesday,
		"Thu": time.Thursday,
		"Fri": time.Friday,
		"Sat": time.Saturday,
		"Sun": time.Sunday,
	}

	dayName := Cast(
		Count(3, IsLetter()),
		func(x []rune) (string, error) {
			return string(x), nil
		},
	)

	return Map(dwDict, dayName)
}

func monthPrefix() p.Combinator[rune, Position, time.Month] {
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

	monthName := Cast(
		Count(3, IsLetter()),
		func(x []rune) (string, error) {
			return string(x), nil
		},
	)

	return Map(monthDict, monthName)
}

func yearWithCentury() p.Combinator[rune, Position, int] {
	return UnsignedN[int](4)
}

func paddedDayNum() p.Combinator[rune, Position, int] {
	return UnsignedN[int](2)
	// return Choice(
	// 	Try(
	// 		Skip(
	// 			OneOf('0', ' '),
	// 			UnsignedN[int](1),
	// 		),
	// 	),
	// 	UnsignedN[int](2),
	// )
}

func paddedHourNum() p.Combinator[rune, Position, int] {
	return UnsignedN[int](2)
}

func paddedMinuteNum() p.Combinator[rune, Position, int] {
	return UnsignedN[int](2)
}

func paddedSecondNum() p.Combinator[rune, Position, int] {
	return UnsignedN[int](2)
}
