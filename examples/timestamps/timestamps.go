package timestamp

import (
	"fmt"
	p "git.sr.ht/~okneniz/parsec/common"
	. "git.sr.ht/~okneniz/parsec/strings"
	t "git.sr.ht/~okneniz/parsec/testing"
	"time"
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

	return Cast(
		Choice(
			Try(SequenceOf('M', 'o', 'n')), // TODO : add String("Mon") helper
			Try(SequenceOf('T', 'u', 'e')),
			Try(SequenceOf('W', 'e', 'd')),
			Try(SequenceOf('T', 'h', 'u')),
			Try(SequenceOf('F', 'r', 'i')),
			Try(SequenceOf('S', 'a', 't')),
			Try(SequenceOf('S', 'u', 'n')),
		),
		func(x []rune) (time.Weekday, error) {
			s := string(x)

			i, exists := dwDict[s]
			if !exists {
				return -1, fmt.Errorf("invalid day of week: %v", s)
			}

			return i, nil
		},
	)
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

	return Cast(
		Choice(
			Try(SequenceOf('J', 'a', 'n')),
			Try(SequenceOf('F', 'e', 'b')),
			Try(SequenceOf('M', 'a', 'r')),
			Try(SequenceOf('A', 'p', 'r')),
			Try(SequenceOf('M', 'a', 'y')),
			Try(SequenceOf('J', 'u', 'n')),
			Try(SequenceOf('J', 'u', 'l')),
			Try(SequenceOf('A', 'u', 'g')),
			Try(SequenceOf('S', 'e', 'p')),
			Try(SequenceOf('O', 'c', 't')),
			Try(SequenceOf('N', 'o', 'v')),
			Try(SequenceOf('D', 'e', 'c')),
		),
		func(x []rune) (time.Month, error) {
			s := string(x)

			i, exists := monthDict[s]
			if !exists {
				return -1, fmt.Errorf("invalid month: %v", s)
			}

			return i, nil
		},
	)
}

func yearWithCentury() p.Combinator[rune, Position, int] {
	return Cast(
		Count(4, IsDigit()),
		t.DigitsToNum,
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
	return Cast(Count(2, IsDigit()), func(x []rune) (int, error) {
		result, err := t.DigitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 23 {
			return -1, fmt.Errorf("invalid hour: %v", result)
		}

		return result, nil
	})
}

func paddedMinuteNum() p.Combinator[rune, Position, int] {
	return Cast(Count(2, IsDigit()), func(x []rune) (int, error) {
		result, err := t.DigitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 59 {
			return -1, fmt.Errorf("invalid minute: %v", result)
		}

		return result, nil
	})
}

func paddedSecondNum() p.Combinator[rune, Position, int] {
	return Cast(Count(2, IsDigit()), func(x []rune) (int, error) {
		result, err := t.DigitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 59 {
			return -1, fmt.Errorf("invalid second: %v", result)
		}

		return result, nil
	})
}

func zoneStr() p.Combinator[rune, Position, *time.Location] {
	return Cast(
		Count(3, Range('A', 'Z')),
		func(x []rune) (*time.Location, error) {
			if string(x) == "LMT" {
				return time.Local, nil
			}

			return time.LoadLocation(string(x))
		},
	)
}

// ANSIC = "Mon Jan _2 15:04:05 2006"
func ansic() p.Combinator[rune, Position, *time.Time] {
	dayOfWeek := dayOfWeekPrefix()
	day := paddedDayNum()
	hour := paddedHourNum()
	separator := Colon()
	minute := paddedMinuteNum()
	second := paddedSecondNum()
	year := yearWithCentury()
	space := IsSpace()

	return func(buffer p.Buffer[rune, Position]) (*time.Time, error) {
		dw, err := dayOfWeek(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		m, err := monthPrefix()(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		d, err := day(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		h, err := hour(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		min, err := minute(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		sec, err := second(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		y, err := year(buffer)
		if err != nil {
			return nil, err
		}

		loc, err := time.LoadLocation("UTC")
		if err != nil {
			return nil, err
		}

		result := time.Date(y, m, d, h, min, sec, 0, loc)
		if result.Weekday() != dw {
			return nil, fmt.Errorf(
				"unexpected day of week: expected %s, actual %v",
				dw,
				result.Weekday(),
			)
		}

		return &result, nil
	}
}

// UnixDate = "Mon Jan _2 15:04:05 MST 2006"
func unixDate() p.Combinator[rune, Position, *time.Time] {
	dayOfWeek := dayOfWeekPrefix()
	space := IsSpace()
	month := monthPrefix()
	day := paddedDayNum()
	hour := paddedHourNum()
	separator := Colon()
	minute := paddedMinuteNum()
	second := paddedSecondNum()
	zone := zoneStr()
	year := yearWithCentury()

	return func(buffer p.Buffer[rune, Position]) (*time.Time, error) {
		dw, err := dayOfWeek(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		m, err := month(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		d, err := day(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		h, err := hour(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		min, err := minute(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		sec, err := second(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		loc, err := zone(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		y, err := year(buffer)
		if err != nil {
			return nil, err
		}

		result := time.Date(y, m, d, h, min, sec, 0, loc)
		if result.Weekday() != dw {
			return nil, fmt.Errorf(
				"unexpected day of week: expected %s, actual %v",
				dw,
				result.Weekday(),
			)
		}

		return &result, nil
	}
}

// RFC1123 = "Mon, 02 Jan 2006 15:04:05 MST"
func rfc1123() p.Combinator[rune, Position, *time.Time] {
	dayOfWeek := dayOfWeekPrefix()
	comma := Comma()
	space := IsSpace()
	day := paddedDayNum()
	month := monthPrefix()
	year := yearWithCentury()
	hour := paddedHourNum()
	separator := Colon()
	minute := paddedMinuteNum()
	second := paddedSecondNum()
	zone := zoneStr()

	return func(buffer p.Buffer[rune, Position]) (*time.Time, error) {
		dw, err := dayOfWeek(buffer)
		if err != nil {
			return nil, err
		}

		_, err = comma(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected comma")
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		d, err := day(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		m, err := month(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		y, err := year(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		h, err := hour(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		min, err := minute(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		sec, err := second(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		loc, err := zone(buffer)
		if err != nil {
			return nil, err
		}

		result := time.Date(y, m, d, h, min, sec, 0, loc)
		if result.Weekday() != dw {
			return nil, fmt.Errorf(
				"unexpected day of week: expected %s, actual %v",
				dw,
				result.Weekday(),
			)
		}

		return &result, nil
	}
}
