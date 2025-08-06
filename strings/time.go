package strings

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/okneniz/parsec/common"
)

// TimeZone - parse one of time zones from passed arguments.
func TimeZone(
	locations ...*time.Location,
) common.Combinator[rune, Position, *time.Location] {
	// TODO : fallback for len=1

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

	return MapStrings(errMessage, m)
}

// TimeZoneByNames - parse one of time zones from passed arguments.
func TimeZoneByNames(
	locationNames ...string,
) (common.Combinator[rune, Position, *time.Location], error) {
	// TODO : fallback for len=1

	locations := make([]*time.Location, 0, len(locationNames))

	for _, locationName := range locationNames {
		loc, err := time.LoadLocation(locationName)
		if err != nil {
			return nil, err
		}

		locations = append(locations, loc)
	}

	return TimeZone(locations...), nil
}
