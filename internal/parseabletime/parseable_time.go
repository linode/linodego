package parseabletime

import (
	"time"
)

const (
	dateLayout = "2006-01-02T15:04:05"
)

type ParseableTime time.Time

func (p *ParseableTime) UnmarshalJSON(b []byte) error {
	var err error
	for _, layout := range []string{time.RFC3339, dateLayout} {
		var t time.Time
		if t, err = time.Parse(`"`+layout+`"`, string(b)); err == nil {
			*p = ParseableTime(t)
			return nil
		}
	}
	return err
}

