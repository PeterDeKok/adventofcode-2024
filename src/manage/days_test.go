package manage

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"testing"
	"time"
)

func TestInitTime(t *testing.T) {
	tz, dt := InitTime("2024-12-01", "Europe/Amsterdam")

	assert.TypeOf[*time.Location](t, tz)
	assert.Equal(t, time.Date(2024, time.December, 1, 5, 0, 0, 0, time.UTC), dt)
}

func TestDaysGenerator_EuropeAmsterdam(t *testing.T) {
	g := DaysGenerator("2024-01-01", "Europe/Amsterdam", 3)

	expected := []string{
		"2024-01-01T06:00:00+01:00",
		"2024-01-02T06:00:00+01:00",
		"2024-01-03T06:00:00+01:00",
	}

	i := 0
	for k, v := range g {
		assert.TypeOf[time.Time](t, v)
		assert.Equal(t, expected[i], v.Format(time.RFC3339))
		assert.Equal(t, i, k)

		i++
	}

	assert.Equal(t, 3, i)
}

func TestDaysGenerator_EST(t *testing.T) {
	g := DaysGenerator("2024-12-25", "EST", 10)

	expected := []string{
		"2024-12-25T00:00:00-05:00",
		"2024-12-26T00:00:00-05:00",
		"2024-12-27T00:00:00-05:00",
		"2024-12-28T00:00:00-05:00",
		"2024-12-29T00:00:00-05:00",
		"2024-12-30T00:00:00-05:00",
		"2024-12-31T00:00:00-05:00",
		"2025-01-01T00:00:00-05:00",
		"2025-01-02T00:00:00-05:00",
		"2025-01-03T00:00:00-05:00",
	}

	i := 0
	for k, v := range g {
		assert.TypeOf[time.Time](t, v)
		assert.Equal(t, expected[i], v.Format(time.RFC3339))
		assert.Equal(t, i, k)

		i++
	}

	assert.Equal(t, 10, i)
}

//func TestInitDays(t *testing.T) {
//	days, lookup := InitPuzzles("2024-12-01", "Europe/Amsterdam", 25)
//
//	assert.TypeOf[[]*Puzzle](t, days)
//	assert.TypeOf[map[string]*Puzzle](t, lookup)
//
//	for i, d := range days {
//		date := d.Day.Format(time.DateOnly)
//
//		t.Run(fmt.Sprintf("slice[%d]:%s", i, date), func(t *testing.T) {
//			// 2024-12-xx 05:00:00 UTC
//			expected := fmt.Sprintf("2024-12-%02d 06:00:00", i+1)
//			actual := d.Day.Format(time.DateTime)
//			if actual != expected {
//				t.Fatalf("slice[%d]: got %v; want %v", i, actual, expected)
//			}
//
//			// 2024-12-xx 05:00:00 UTC
//			expectedTz := "Europe/Amsterdam"
//			actualTz := d.Day.Location().String()
//			if actualTz != expectedTz {
//				t.Fatalf("slice[%d]: got location %v; want %v", i, actualTz, expectedTz)
//			}
//
//			dateOnly := d.Day.Format(time.DateOnly)
//			if _, ok := lookup[dateOnly]; !ok {
//				t.Fatalf("slice[%d]: %s not found in lookup %v", i, dateOnly, expectedTz)
//			}
//
//			assert.TypeOf[*Part](t, d.Part1)
//			assert.NotNil(t, d.Part1)
//			assert.TypeOf[*Part](t, d.Part2)
//			assert.NotNil(t, d.Part2)
//		})
//	}
//
//	assert.HasLen(t, 25, days)
//}
