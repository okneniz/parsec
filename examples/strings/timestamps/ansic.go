package timestamp

import (
	"fmt"
	"time"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/strings"
)

// ANSIC = "Mon Jan _2 15:04:05 2006"
func ansic() parsec.Combinator[rune, strings.Position, *time.Time, parsec.Stateless] {
	dayOfWeek := dayOfWeekPrefix()
	day := paddedDayNum()
	hour := paddedHourNum()
	separator := strings.Colon[parsec.Stateless]()
	minute := paddedMinuteNum()
	second := paddedSecondNum()
	year := yearWithCentury()

	return func(
		state parsec.Stateless,
		buffer parsec.Buffer[rune, strings.Position],
	) (*time.Time, parsec.Error[strings.Position]) {
		dw, err := dayOfWeek(state, buffer)
		if err != nil {
			return nil, err
		}

		m, err := monthPrefix()(state, buffer)
		if err != nil {
			return nil, err
		}

		d, err := day(state, buffer)
		if err != nil {
			return nil, err
		}

		h, err := hour(state, buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(state, buffer)
		if err != nil {
			return nil, err
		}

		min, err := minute(state, buffer)
		if err != nil {
			return nil, err
		}

		_, err = separator(state, buffer)
		if err != nil {
			return nil, err
		}

		sec, err := second(state, buffer)
		if err != nil {
			return nil, err
		}

		y, err := year(state, buffer)
		if err != nil {
			return nil, err
		}

		loc, sysErr := time.LoadLocation("UTC")
		if sysErr != nil {
			return nil, parsec.NewParseError(buffer.Position(), sysErr.Error())
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
