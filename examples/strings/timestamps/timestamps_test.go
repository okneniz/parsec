package timestamp

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	. "github.com/okneniz/parsec/strings"
	. "github.com/okneniz/parsec/testing"
)

func TestTimestamps(t *testing.T) {
	seed := time.Now().UnixNano()
	t.Log("seed: ", seed)

	source := rand.New(rand.NewSource(seed))
	r := rand.New(source)

	dates := randomDates(r, 1000)
	formattedDates := randomFormattedDates(r, dates)
	input := JoinBy(func() string { return "\n" }, formattedDates...)

	t.Log("input:")
	t.Logf("%#v", input)

	oneOfDate := Choice(
		Try(ansic()),
		Try(unixDate()),
		Try(rfc1123()),
	)

	comb := SepBy(len(dates), oneOfDate, Eq('\n'))

	result, err := ParseString(input, comb)
	Check(t, err)
	AssertEqDump(t, result, dates)
}

func BenchmarkNativeUnixDate(b *testing.B) {
	seed := time.Now().UnixNano()
	source := rand.New(rand.NewSource(seed))
	r := rand.New(source)

	b.Log("seed: ", seed)
	layout := "Mon Jan _2 15:04:05 MST 2006"

	dates := make([]string, b.N)
	for i, d := range randomDates(r, b.N) {
		dates[i] = d.Format(layout)
	}

	b.ResetTimer()
	for _, input := range dates {
		time.Parse(layout, input)
	}
}

func BenchmarkParsecUnixDate(b *testing.B) {
	seed := time.Now().UnixNano()
	source := rand.New(rand.NewSource(seed))
	r := rand.New(source)

	b.Log("seed: ", seed)
	layout := "Mon Jan _2 15:04:05 MST 2006"

	dates := make([]string, b.N)
	for i, d := range randomDates(r, b.N) {
		dates[i] = d.Format(layout)
	}

	comb := unixDate()

	b.ResetTimer()
	for _, input := range dates {
		ParseString(input, comb)
	}
}

func randomFormattedDates(r *rand.Rand, dt []*time.Time) []string {
	result := make([]string, len(dt))
	for i, d := range dt {
		l := randomLayout(r)
		result[i] = d.Format(l)
	}
	return result
}

func randomDates(r *rand.Rand, count int) []*time.Time {
	result := make([]*time.Time, count)
	for i := 0; i < count; i++ {
		result[i] = randomDate(r)
	}
	return result
}

func randomDate(r *rand.Rand) *time.Time {
	year := randomInt(r, 1000, 2025)
	month := time.Month(randomInt(r, 1, 12))
	day := randomInt(r, 10, 28)
	hour := randomInt(r, 0, 23)
	min := randomInt(r, 0, 59)
	sec := randomInt(r, 0, 59)
	loc := randomLocation(r)
	d := time.Date(year, month, day, hour, min, sec, 0, loc)
	return &d
}

var (
	allowZones = []string{
		"UTC", "EST", "GMT",
	}

	allowLayouts = []string{
		// only layouts with time zone to avoid losing data
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}
)

func randomLocation(r *rand.Rand) *time.Location {
	x := randomInt(r, 0, len(allowZones)-1)
	loc, err := time.LoadLocation(allowZones[x])
	if err != nil {
		panic(fmt.Sprintf("generate random location error: %s (x=%v,name=%v)", err, x, allowZones[x]))
	}
	return loc
}

func randomInt(r *rand.Rand, from, to int) int {
	return from + r.Intn(to-from+1)
}

func randomLayout(r *rand.Rand) string {
	x := randomInt(r, 0, len(allowLayouts)-1)
	return allowLayouts[x]
}
