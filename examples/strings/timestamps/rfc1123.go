package timestamp

import (
	"fmt"
	"time"

	"github.com/okneniz/parsec/common"
	"github.com/okneniz/parsec/strings"
)

// RFC1123 = "Mon, 02 Jan 2006 15:04:05 MST"
func rfc1123() common.Combinator[rune, strings.Position, *time.Time] {
	dayOfWeek := dayOfWeekPrefix()
	comma := strings.Comma()
	space := strings.Space("expected space")
	day := paddedDayNum()
	month := monthPrefix()
	year := yearWithCentury()
	hour := paddedHourNum()
	separator := strings.Colon()
	minute := paddedMinuteNum()
	second := paddedSecondNum()
	zone, _ := strings.TimeZoneByNames("UTC", "EST", "GMT")

	return func(
		buffer common.Buffer[rune, strings.Position],
	) (*time.Time, common.Error[strings.Position]) {
		dw, err := dayOfWeek(buffer)
		if err != nil {
			return nil, err
		}

		_, err = comma(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, err
		}

		d, err := day(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, err
		}

		m, err := month(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
		if err != nil {
			return nil, err
		}

		y, err := year(buffer)
		if err != nil {
			return nil, err
		}

		_, err = space(buffer)
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

		_, err = space(buffer)
		if err != nil {
			return nil, err
		}

		loc, err := zone(buffer)
		if err != nil {
			return nil, err
		}

		result := time.Date(y, m, d, h, min, sec, 0, loc)
		if result.Weekday() != dw {
			return nil, common.NewParseError(
				buffer.Position(),
				fmt.Sprintf("unexpected day of week: expected %s, actual %v", dw, result.Weekday()),
			)
		}

		return &result, nil
	}
}
