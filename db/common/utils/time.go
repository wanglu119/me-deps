package utils

import (
	"fmt"
	"time"
)

type Time time.Time

const (
	timeFormartWithTime = "2006-01-02 15:04:05"
	timeFormartNoTime   = "2006-01-02"
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// +2 为双引号
	if len(data) == 19+2 {
		now, err := time.ParseInLocation(`"`+timeFormartWithTime+`"`, string(data), time.Local)
		*t = Time(now)
		return err
	}

	if len(data) == 10+2 {
		now, err := time.ParseInLocation(`"`+timeFormartNoTime+`"`, string(data), time.Local)
		*t = Time(now)
		return err
	}

	err = fmt.Errorf("unsupport date format: %s, length: %d", string(data), len(data))

	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormartWithTime)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormartWithTime)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormartWithTime)
}

func (t Time) Format(format string) string {
	return time.Time(t).Format(format)
}

func (t Time) Time() time.Time {
	return time.Time(t)
}
