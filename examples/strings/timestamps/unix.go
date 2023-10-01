package timestamp

import (
	"fmt"
	"time"

	p "github.com/okneniz/parsec/common"
	. "github.com/okneniz/parsec/strings"
)

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
	year := yearWithCentury()
	zone, _ := TimeZoneByNames("UTC", "EST", "GMT")

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
