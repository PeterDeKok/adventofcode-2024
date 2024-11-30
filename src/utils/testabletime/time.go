package testabletime

import "time"

var testDiff time.Duration = 0

func SetTimeOffset(diff time.Duration) {
	testDiff = diff
}

func Now() time.Time {
	return time.Now().Add(testDiff)
}
