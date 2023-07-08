package strings

import (
	"time"

	p "git.sr.ht/~okneniz/parsec/common"
)

func TimeZone(
	locations ...*time.Location,
) p.Combinator[rune, Position, *time.Location] {
	m := make(map[string]*time.Location, len(locations))
	names := make([]string, len(locations))

	t := time.Now()

	for i, loc := range locations {
		tt := t.In(loc)
		zoneName, _ := tt.Zone()
		m[zoneName] = loc
		names[i] = zoneName
	}

	return Map(m, OneOfStrings(names...))
}

func TimeZoneByNames(
	locationNames ...string,
) (p.Combinator[rune, Position, *time.Location], error) {
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
