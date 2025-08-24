package strings

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/stretchr/testify/assert"
)

func TestTimeZone(t *testing.T) {
	zones := []string{"UTC", "EST", "GMT"}

	t.Run("case 1", func(t *testing.T) {
		comb, err := TimeZoneByNames(zones...)
		if err != nil {
			t.Fatal(err)
		}

		for _, zone := range zones {
			result, err := ParseString(zone, comb)
			assert.NoError(t, err)
			assert.NotNil(t, result, "expected pointer to time zone")
			assert.Equal(t, result.String(), zone)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		comb, err := TimeZoneByNames("Europe/Moscow")
		if err != nil {
			t.Fatal(err)
		}

		for _, zone := range zones {
			result, err := ParseString(zone, comb)
			assert.Error(t, err)
			assert.Nil(t, result, nil)
		}
	})
}

func TestTimeZoneByNames(t *testing.T) {
	zones := []string{"UTC", "EST", "GMT"}

	locationsM := make(map[string]*time.Location, len(zones))
	locations := make([]*time.Location, 0, len(zones))

	for _, zone := range zones {
		loc, err := time.LoadLocation(zone)
		if err != nil {
			t.Fatal(err)
		}

		locationsM[zone] = loc
		locations = append(locations, loc)
	}

	t.Run("case 1", func(t *testing.T) {
		comb := TimeZone(locations...)

		for _, zone := range zones {
			result, err := ParseString(zone, comb)
			assert.NoError(t, err)
			assert.NotNil(t, result, "expected pointer to time zone")
			assert.Equal(t, result.String(), zone)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			t.Fatal(err)
		}

		comb := TimeZone(loc)

		for _, zone := range zones {
			result, err := ParseString(zone, comb)
			assert.Error(t, err)
			assert.Nil(t, result)
		}
	})
}
