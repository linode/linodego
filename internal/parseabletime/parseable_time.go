package parseabletime

import (
	"encoding/json"
	"time"
)

const (
	dateLayout = "2006-01-02T15:04:05"
)

type ParseableTime time.Time

func (p *ParseableTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(`"`+dateLayout+`"`, string(b))
	if err != nil {
		return err
	}

	*p = ParseableTime(t)

	return nil
}

func (p *ParseableTime) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}
	t := time.Time(*p)
	return json.Marshal(t.Format(dateLayout))
}
