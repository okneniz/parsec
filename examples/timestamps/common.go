package timestamp

import (
	"fmt"
	"time"

	p "git.sr.ht/~okneniz/parsec/common"
	. "git.sr.ht/~okneniz/parsec/strings"
	t "git.sr.ht/~okneniz/parsec/testing"
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
	return Cast(
		UnsignedN[uint](4),
		func(x uint) (int, error) {
			return int(x), nil
		},
	)
}

func paddedDayNum() p.Combinator[rune, Position, int] {
	pad := OneOf('0', ' ')

	return Cast(
		Choice(
			Count(2, IsDigit()),
			Skip(Optional(pad, rune(0)), Some(1, IsDigit())),
		),
		func(x []rune) (int, error) {
			result, err := t.DigitsToNum(x)
			if err != nil {
				return -1, err
			}
			if result <= 0 || result > 31 {
				return -1, fmt.Errorf("invalid day num: %v (%v)", result, string(x))
			}

			return result, nil
		},
	)
}

func paddedHourNum() p.Combinator[rune, Position, int] {
	return Cast(
		UnsignedN[uint](2),
		func(x uint) (int, error) {
			return int(x), nil
		},
	)
}

func paddedMinuteNum() p.Combinator[rune, Position, int] {
	return Cast(
		UnsignedN[uint](2),
		func(x uint) (int, error) {
			return int(x), nil
		},
	)
}

func paddedSecondNum() p.Combinator[rune, Position, int] {
	return Cast(
		UnsignedN[uint](2),
		func(x uint) (int, error) {
			return int(x), nil
		},
	)
}

// TODO : add helper with available zones from file
func zoneStr() p.Combinator[rune, Position, *time.Location] {
	return Cast(
		Count(3, Range('A', 'Z')),
		func(x []rune) (*time.Location, error) {
			return time.LoadLocation(string(x))
		},
	)
}
