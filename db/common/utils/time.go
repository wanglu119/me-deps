package utils

import (
	"time"
)

type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) > 2 {
		now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
		*t = Time(now)
		return err
	}

	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormart)
}

func (t Time) Format(format string) string {
	return time.Time(t).Format(format)
}

func (t Time) Time() time.Time {
	return time.Time(t)
}
