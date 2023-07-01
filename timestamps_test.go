package parsec

import (
	"fmt"
	"time"
	"math/rand"
	"math"
	"testing"
)

func TestTimestamps(t *testing.T) {
	digit := Range(true, byte('0'), byte('9'))

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
			Try(Sequence(3, Eq(true, byte('M')), Eq(true, byte('o')), Eq(true, byte('n')))),
			Try(Sequence(3, Eq(true, byte('T')), Eq(true, byte('u')), Eq(true, byte('e')))),
			Try(Sequence(3, Eq(true, byte('W')), Eq(true, byte('e')), Eq(true, byte('d')))),
			Try(Sequence(3, Eq(true, byte('T')), Eq(true, byte('h')), Eq(true, byte('u')))),
			Try(Sequence(3, Eq(true, byte('F')), Eq(true, byte('r')), Eq(true, byte('i')))),
			Try(Sequence(3, Eq(true, byte('S')), Eq(true, byte('a')), Eq(true, byte('t')))),
			Try(Sequence(3, Eq(true, byte('S')), Eq(true, byte('u')), Eq(true, byte('n')))),
		),
		func(x []byte) (time.Weekday, error) {
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
			Try(Sequence(3, Eq(true, byte('J')), Eq(true, byte('a')), Eq(true, byte('n')))),
			Try(Sequence(3, Eq(true, byte('F')), Eq(true, byte('e')), Eq(true, byte('b')))),
			Try(Sequence(3, Eq(true, byte('M')), Eq(true, byte('a')), Eq(true, byte('r')))),
			Try(Sequence(3, Eq(true, byte('A')), Eq(true, byte('p')), Eq(true, byte('r')))),
			Try(Sequence(3, Eq(true, byte('M')), Eq(true, byte('a')), Eq(true, byte('y')))),
			Try(Sequence(3, Eq(true, byte('J')), Eq(true, byte('u')), Eq(true, byte('n')))),
			Try(Sequence(3, Eq(true, byte('J')), Eq(true, byte('u')), Eq(true, byte('l')))),
			Try(Sequence(3, Eq(true, byte('A')), Eq(true, byte('u')), Eq(true, byte('g')))),
			Try(Sequence(3, Eq(true, byte('S')), Eq(true, byte('e')), Eq(true, byte('p')))),
			Try(Sequence(3, Eq(true, byte('O')), Eq(true, byte('c')), Eq(true, byte('t')))),
			Try(Sequence(3, Eq(true, byte('N')), Eq(true, byte('o')), Eq(true, byte('v')))),
			Try(Sequence(3, Eq(true, byte('D')), Eq(true, byte('e')), Eq(true, byte('c')))),
		),
		func(x []byte) (time.Month, error) {
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
		digitsToNum,
	))

	pad := OneOf(true, byte('0'), byte(' '))

	paddedDayNum := Trace(t, "day",
		Cast(
			Choice(
				Count(2, digit),
				Skip(Optional(pad, byte(0)), Some(1, digit)),
			),
			func(x []byte) (int, error) {
				result, err := digitsToNum(x)
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

	paddedHourNum := Trace(t, "hour", Cast(Count(2, digit), func(x []byte) (int, error) {
		result, err := digitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 23 {
			return -1, fmt.Errorf("invalid hour: %v", result)
		}

		return result, nil
	}))

	paddedMinuteNum := Trace(t, "minute", Cast(Count(2, digit), func(x []byte) (int, error) {
		result, err := digitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 59 {
			return -1, fmt.Errorf("invalid minute: %v", result)
		}

		return result, nil
	}))

	paddedSecondNum := Trace(t, "second", Cast(Count(2, digit), func(x []byte) (int, error) {
		result, err := digitsToNum(x)
		if err != nil {
			return -1, err
		}
		if result < 0 || result > 59 {
			return -1, fmt.Errorf("invalid second: %v", result)
		}

		return result, nil
	}))

	zoneStr := Trace(t, "zone", Cast(
		Count(3, Range(true, byte('A'), byte('Z'))),
		func(x []byte) (*time.Location, error) {
			return time.LoadLocation(string(x))
		},
	))

	space := Eq(true, byte(' '))
	sep := Eq(true, byte(':'))
	comma := Eq(true, byte(','))

	// ANSIC = "Mon Jan _2 15:04:05 2006"
	ansic := Trace(t, "ansic", func(buffer Buffer[byte]) (*time.Time, error) {
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
	unixDate := Trace(t, "unix date", func(buffer Buffer[byte]) (*time.Time, error) {
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
	rfc1123 := Trace(t, "RFC1123", func(buffer Buffer[byte]) (*time.Time, error) {
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

	dates := randomDates(r, 300)
	formattedDates := randomFormattedDates(r, dates)
	input := joinBy(func() string { return "\n" }, formattedDates...)

	t.Log("input:")
	t.Logf("%#v", input)

	oneOfDate := Choice(
		Try(ansic),
		Try(unixDate),
		Try(rfc1123),
	)

	comb := SepBy(len(dates), oneOfDate, Eq(true, byte('\n')))

	result, err := ParseBytes([]byte(input), comb)
	check(t, err)
	assertEqDump(t, result, dates)
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
		"UTC",
		// "Europe/London",
		// "MST",
	}

	allowLayouts = []string{
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}
)

func randomLocation(r *rand.Rand) *time.Location {
	x := randomInt(r, 0, len(allowZones) - 1)
	loc, err := time.LoadLocation(allowZones[x])
	if err != nil {
		panic(fmt.Sprintf("generate random location error: %s (x=%v,name=%v)", err, x, allowZones[x]))
	}
	return loc
}

func randomInt(r *rand.Rand, from, to int) int {
	return from + r.Intn(to - from)
}

func randomLayout(r *rand.Rand) string {
	x := randomInt(r, 0, len(allowLayouts) - 1)
	return allowLayouts[x]
}

func digitsToNum(ds []byte) (int, error) {
	if len(ds) == 0 {
		return -1, fmt.Errorf("invalid number []bytes: %v, string: %v", ds, string(ds))
	}

	m := map[byte]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
	}

	num := m[ds[len(ds)-1]]

	for i, d := range ds[:len(ds)-1] {
		l := math.Pow(10, float64(len(ds) - i - 1))
		v := int(l) * m[d]
		num += v
	}

	return num, nil
}


func TestDigitsToNum(t *testing.T) {
	examples := map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"10": 10,
		"11": 11,
		"101": 101,
		"123": 123,
		"10723": 10723,
		"50221": 50221,
	}

	for input, expected := range examples {
		actual, err := digitsToNum([]byte(input))
		if err != nil {
			t.Error(err)
			t.Errorf("expected %v, actual %v - input %v", expected, actual, input)
			continue
		}
		if expected != actual {
			t.Errorf("expected %v, actual %v - input %v", expected, actual, input)
		}
	}
}
