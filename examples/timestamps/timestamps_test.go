package timestamp

import (
	"fmt"
	. "git.sr.ht/~okneniz/parsec/strings"
	p "git.sr.ht/~okneniz/parsec/common"
	. "git.sr.ht/~okneniz/parsec/testing"
	"testing"
	"time"
	"math/rand"
)

func TestTimestamps(t *testing.T) {
	digit := Range('0', '9')

	dwDict := map[string]time.Weekday{
		"Mon": time.Monday,
		"Tue": time.Tuesday,
		"Wed": time.Wednesday,
		"Thu": time.Thursday,
		"Fri": time.Friday,
		"Sat": time.Saturday,
		"Sun": time.Sunday,
	}

	dayOfWeekPrefix := Trace(t, "day of week", Cast(
		Choice(
			Try(SequenceOf('M', 'o', 'n')),
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
	))

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

	monthPrefix := Trace(t, "month prefix", Cast(
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
	))

	yearWithCentury := Trace(t, "year with century", Cast(
		Count(4, digit),
		DigitsToNum,
	))

	pad := OneOf('0', ' ')

	paddedDayNum := Trace(t, "day",
		Cast(
			Choice(
				Count(2, digit),
				Skip(Optional(pad, rune(0)), Some(1, digit)),
			),
			func(x []rune) (int, error) {
				result, err := DigitsToNum(x)
				if err != nil {
					return -1, err
				}
				if result <= 0 || result > 31 {
					return -1, fmt.Errorf("invalid day num: %v (%v)", result, string(x))
				}

				return result, nil
			},
		),
	)

	paddedHourNum := Trace(t, "hour", Cast(Count(2, digit), func(x []rune) (int, error) {
		result, err := DigitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 23 {
			return -1, fmt.Errorf("invalid hour: %v", result)
		}

		return result, nil
	}))

	paddedMinuteNum := Trace(t, "minute", Cast(Count(2, digit), func(x []rune) (int, error) {
		result, err := DigitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 59 {
			return -1, fmt.Errorf("invalid minute: %v", result)
		}

		return result, nil
	}))

	paddedSecondNum := Trace(t, "second", Cast(Count(2, digit), func(x []rune) (int, error) {
		result, err := DigitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 59 {
			return -1, fmt.Errorf("invalid second: %v", result)
		}

		return result, nil
	}))

	zoneStr := Trace(t, "zone", Cast(
		Count(3, Range('A', 'Z')),
		func(x []rune) (*time.Location, error) {
			if string(x) == "LMT" {
				return time.Local, nil
			}

			return time.LoadLocation(string(x))
		},
	))

	space := Eq(' ')
	sep := Eq(':')
	comma := Eq(',')

	// ANSIC = "Mon Jan _2 15:04:05 2006"
	ansic := Trace(t, "ansic", func(buffer p.Buffer[rune, Position]) (*time.Time, error) {
		dw, err := dayOfWeekPrefix(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		m, err := monthPrefix(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		d, err := paddedDayNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		h, err := paddedHourNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = sep(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		min, err := paddedMinuteNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = sep(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		sec, err := paddedSecondNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		y, err := yearWithCentury(buffer)
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
	})

	// UnixDate = "Mon Jan _2 15:04:05 MST 2006"
	unixDate := Trace(t, "unix date", func(buffer p.Buffer[rune, Position]) (*time.Time, error) {
		dw, err := dayOfWeekPrefix(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		m, err := monthPrefix(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		d, err := paddedDayNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		h, err := paddedHourNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = sep(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		min, err := paddedMinuteNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = sep(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		sec, err := paddedSecondNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		loc, err := zoneStr(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		y, err := yearWithCentury(buffer)
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
	})

	// RFC1123 = "Mon, 02 Jan 2006 15:04:05 MST"
	rfc1123 := Trace(t, "RFC1123", func(buffer p.Buffer[rune, Position]) (*time.Time, error) {
		dw, err := dayOfWeekPrefix(buffer)
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

		d, err := paddedDayNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		m, err := monthPrefix(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		y, err := yearWithCentury(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		h, err := paddedHourNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = sep(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		min, err := paddedMinuteNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = sep(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected separator ':'")
		}

		sec, err := paddedSecondNum(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, fmt.Errorf("expected space")
		}

		loc, err := zoneStr(buffer)
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
	})

	seed := time.Now().UnixNano()
	t.Log("seed: ", seed)

	source := rand.New(rand.NewSource(seed))
	r := rand.New(source)

	dates := randomDates(r, 1000)
	formattedDates := randomFormattedDates(r, dates)
	input := JoinBy(func() string { return "\n" }, formattedDates...)

	t.Log("input:")
	t.Logf("%#v", input)

	oneOfDate := Choice(
		Try(ansic),
		Try(unixDate),
		Try(rfc1123),
	)

	comb := SepBy(len(dates), oneOfDate, Eq('\n'))

	result, err := ParseString(input, comb)
	Check(t, err)
	AssertEqDump(t, result, dates)
}

func randomFormattedDates(r *rand.Rand, dt []*time.Time) []string {
	result := make([]string, len(dt))
	for i, d := range dt {
		l := randomLayout(r)
		result[i] = d.Format(l)
	}
	return result
}

func randomDates(r *rand.Rand, count int) []*time.Time {
	result := make([]*time.Time, count)
	for i := 0; i < count; i++ {
		result[i] = randomDate(r)
	}
	return result
}

func randomDate(r *rand.Rand) *time.Time {
	year := randomInt(r, 1000, 2025)
	month := time.Month(randomInt(r, 1, 12))
	day := randomInt(r, 10, 28)
	hour := randomInt(r, 0, 23)
	min := randomInt(r, 0, 59)
	sec := randomInt(r, 0, 59)
	loc := randomLocation(r)
	d := time.Date(year, month, day, hour, min, sec, 0, loc)
	return &d
}

var (
	allowZones = []string{
		"UTC",
		"MST",
	}

	allowLayouts = []string{
		// "Mon Jan _2 15:04:05 2006", TODO : only layouts with time zone to avoid losing data
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}
)

func randomLocation(r *rand.Rand) *time.Location {
	x := randomInt(r, 0, len(allowZones)-1)
	loc, err := time.LoadLocation(allowZones[x])
	if err != nil {
		panic(fmt.Sprintf("generate random location error: %s (x=%v,name=%v)", err, x, allowZones[x]))
	}
	return loc
}

func randomInt(r *rand.Rand, from, to int) int {
	return from + r.Intn(to-from)
}

func randomLayout(r *rand.Rand) string {
	x := randomInt(r, 0, len(allowLayouts)-1)
	return allowLayouts[x]
}
