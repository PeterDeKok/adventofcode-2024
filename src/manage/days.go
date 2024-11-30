package manage

import (
	"fmt"
	"iter"
	"time"
)

func InitTime(firstday, tz string) (*time.Location, time.Time) {
	europeAmsterdam, err := time.LoadLocation(tz)
	if err != nil {
		panic(fmt.Errorf("failed to init time: %v", err))
	}

	dayCursor, err := time.ParseInLocation(time.DateTime, firstday+" 05:00:00", time.UTC)
	if err != nil {
		panic(fmt.Errorf("failed to init time: %v", err))
	}

	return europeAmsterdam, dayCursor
}

func DaysGenerator(firstday, tz string, nrDays int) iter.Seq2[int, time.Time] {
	return func(yield func(int, time.Time) bool) {
		loc, day := InitTime(firstday, tz)

		for i := 0; i < nrDays; i++ {
			yield(i, day.In(loc))

			day = day.AddDate(0, 0, 1)
		}
	}
}
