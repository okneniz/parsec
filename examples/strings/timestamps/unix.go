package timestamp

import (
	"fmt"
	"time"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// UnixDate = "Mon Jan _2 15:04:05 MST 2006"
func unixDate() parsec.Combinator[rune, strings.Position, *time.Time] {
	dayOfWeek := dayOfWeekPrefix()
	month := monthPrefix()
	day := paddedDayNum()
	hour := paddedHourNum()
	separator := strings.Colon()
	minute := paddedMinuteNum()
	second := paddedSecondNum()
	year := yearWithCentury()
	zone, _ := strings.TimeZoneByNames("UTC", "EST", "GMT")

	return func(
		buffer parsec.Buffer[rune, strings.Position],
	) (*time.Time, parsec.Error[strings.Position]) {
		dw, err := dayOfWeek(buffer)
		if err != nil {
			return nil, err
		}

		m, err := month(buffer)
		if err != nil {
			return nil, err
		}

		d, err := day(buffer)
		if err != nil {
			return nil, err
		}

		h, err := hour(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, err
		}

		min, err := minute(buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(buffer)
		if err != nil {
			return nil, err
		}

		sec, err := second(buffer)
		if err != nil {
			return nil, err
		}

		loc, err := zone(buffer)
		if err != nil {
			return nil, err
		}

		y, err := year(buffer)
		if err != nil {
			return nil, err
		}

		result := time.Date(y, m, d, h, min, sec, 0, loc)
		if result.Weekday() != dw {
			return nil, parsec.NewParseError(
				buffer.Position(),
				fmt.Sprintf(
					"unexpected day of week: expected %s, actual %v",
					dw,
					result.Weekday(),
				),
			)
		}

		return &result, nil
	}
}
