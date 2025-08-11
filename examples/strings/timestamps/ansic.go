package timestamp

import (
	"fmt"
	"time"

	"github.com/okneniz/parsec/common"
	"github.com/okneniz/parsec/strings"
)

// ANSIC = "Mon Jan _2 15:04:05 2006"
func ansic() common.Combinator[rune, strings.Position, *time.Time] {
	dayOfWeek := dayOfWeekPrefix()
	day := paddedDayNum()
	hour := paddedHourNum()
	separator := strings.Colon()
	minute := paddedMinuteNum()
	second := paddedSecondNum()
	year := yearWithCentury()

	return func(
		buffer common.Buffer[rune, strings.Position],
	) (*time.Time, common.Error[strings.Position]) {
		dw, err := dayOfWeek(buffer)
		if err != nil {
			return nil, err
		}

		m, err := monthPrefix()(buffer)
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

		y, err := year(buffer)
		if err != nil {
			return nil, err
		}

		loc, sysErr := time.LoadLocation("UTC")
		if sysErr != nil {
			return nil, common.NewParseError(buffer.Position(), sysErr.Error())
		}

		result := time.Date(y, m, d, h, min, sec, 0, loc)
		if result.Weekday() != dw {
			return nil, common.NewParseError(
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
