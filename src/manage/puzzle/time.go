package puzzle

import (
	"encoding/json"
	"time"
)

type Date time.Time
type Timestamp time.Time

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ts).Unix())
}

func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	var v int64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*ts = Timestamp(time.Unix(v, 0))

	return nil
}

func (ts *Timestamp) Before(other *Timestamp) bool {
	if ts == nil || other == nil {
		panic("No nil check for timestamp values")
	}

	return time.Time(*ts).Before(time.Time(*other))
}

func (ts *Timestamp) IsUnix() bool {
	if ts == nil {
		return false
	}

	return time.Time(*ts).Equal(time.Unix(0, 0))
}

func (ts *Timestamp) Format(layout string) string {
	if ts == nil {
		return ""
	}

	return time.Time(*ts).Format(layout)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format(time.DateOnly))
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var v string
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	t, err := time.Parse(time.DateOnly, v)
	if err != nil {
		return err
	}

	*d = Date(t)

	return nil
}
