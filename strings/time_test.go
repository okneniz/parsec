package strings

import (
	"testing"
	"time"

	"github.com/okneniz/parsec"

	_ "time/tzdata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeZone(t *testing.T) {
	zones := []string{"UTC", "EST", "GMT"}

	t.Run("case 1", func(t *testing.T) {
		comb, err := TimeZoneByNames[parsec.Stateless](zones...)
		if err != nil {
			t.Fatal(err)
		}

		for _, zone := range zones {
			result, err := ParseString(parsec.Stateless{}, zone, comb)
			require.NoError(t, err)
			assert.NotNil(t, result, "expected pointer to time zone")
			assert.Equal(t, result.String(), zone)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		comb, err := TimeZoneByNames[parsec.Stateless]("Europe/Moscow")
		if err != nil {
			t.Fatal(err)
		}

		for _, zone := range zones {
			result, err := ParseString(parsec.Stateless{}, zone, comb)
			require.Error(t, err)
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
		comb := TimeZone[parsec.Stateless](locations...)

		for _, zone := range zones {
			result, err := ParseString(parsec.Stateless{}, zone, comb)
			require.NoError(t, err)
			assert.NotNil(t, result, "expected pointer to time zone")
			assert.Equal(t, result.String(), zone)
		}
	})

	t.Run("case 2", func(t *testing.T) {
		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			t.Fatal(err)
		}

		comb := TimeZone[parsec.Stateless](loc)

		for _, zone := range zones {
			result, err := ParseString(parsec.Stateless{}, zone, comb)
			require.Error(t, err)
			assert.Nil(t, result)
		}
	})
}
