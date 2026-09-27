package strings

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/okneniz/parsec"
)

// TimeZone parses the current zone name of one of the passed locations,
// for example "UTC" or "MSK", and returns the matching location.
func TimeZone[S any](
	locations ...*time.Location,
) parsec.Combinator[rune, Position, *time.Location, S] {
	m := make(map[string]*time.Location, len(locations))
	names := make([]string, len(locations))

	t := time.Now()

	for i, loc := range locations {
		tt := t.In(loc)
		zoneName, _ := tt.Zone()
		m[zoneName] = loc
		names[i] = zoneName
	}

	sort.SliceStable(
		names,
		func(i, j int) bool { return names[i] < names[j] },
	)

	errMessage := fmt.Sprintf(
		"expected one of time zones: %s",
		strings.Join(names, ","),
	)

	return MapStrings[*time.Location, S](errMessage, m)
}

// TimeZoneByNames loads the locations by their IANA names
// (for example "UTC" or "Europe/Moscow") and returns a combinator
// parsing the current zone name of one of them.
func TimeZoneByNames[S any](
	locationNames ...string,
) (parsec.Combinator[rune, Position, *time.Location, S], error) {
	locations := make([]*time.Location, 0, len(locationNames))

	for _, locationName := range locationNames {
		loc, err := time.LoadLocation(locationName)
		if err != nil {
			return nil, err
		}

		locations = append(locations, loc)
	}

	return TimeZone[S](locations...), nil
}
