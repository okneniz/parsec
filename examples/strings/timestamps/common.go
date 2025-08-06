package timestamp

import (
	"time"

	"github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
)

func dayOfWeekPrefix() common.Combinator[rune, Position, time.Weekday] {
	dwDict := map[string]time.Weekday{
		"Mon": time.Monday,
		"Tue": time.Tuesday,
		"Wed": time.Wednesday,
		"Thu": time.Thursday,
		"Fri": time.Friday,
		"Sat": time.Saturday,
		"Sun": time.Sunday,
	}

	return MapStrings("expected day of week", dwDict)
}

func monthPrefix() common.Combinator[rune, Position, time.Month] {
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

	return MapStrings("expected name of month", monthDict)
}

func yearWithCentury() common.Combinator[rune, Position, int] {
	return UnsignedN[int](4, "expected year with century")
}

func paddedDayNum() common.Combinator[rune, Position, int] {
	return UnsignedN[int](2, "expected day number")
}

func paddedHourNum() common.Combinator[rune, Position, int] {
	return UnsignedN[int](2, "expected hour number")
}

func paddedMinuteNum() common.Combinator[rune, Position, int] {
	return UnsignedN[int](2, "expected minute number")
}

func paddedSecondNum() common.Combinator[rune, Position, int] {
	return UnsignedN[int](2, "expected second number")
}
