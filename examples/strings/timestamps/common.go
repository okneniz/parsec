package timestamp

import (
	"time"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
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

	return MapStrings(dwDict)
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

	return MapStrings(monthDict)
}

func yearWithCentury() p.Combinator[rune, Position, int] {
	return UnsignedN[int](4)
}

func paddedDayNum() p.Combinator[rune, Position, int] {
	return UnsignedN[int](2)
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
