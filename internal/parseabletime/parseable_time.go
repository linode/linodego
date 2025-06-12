package parseabletime

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	dateLayout = "2006-01-02T15:04:05"
)

type ParseableTime time.Time

var ErrParseableTime = errors.New("parseable time error: invalid format, expected ISO8601")

func (p *ParseableTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return errors.Join(ErrParseableTime, err)
	}

	t, err := time.Parse(dateLayout, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return errors.Join(ErrParseableTime, err)
		}
	}
	*p = ParseableTime(t)

	return nil
}
