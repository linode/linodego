package duration

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalTimeRemaining(t *testing.T) {
	if *UnmarshalTimeRemaining(json.RawMessage("\"1:23\"")) != 83 {
		t.Errorf("Error parsing duration style time_remaining")
	}
	if UnmarshalTimeRemaining(json.RawMessage("null")) != nil {
		t.Errorf("Error parsing null time_remaining")
	}
	if *UnmarshalTimeRemaining(json.RawMessage("0")) != 0 {
		t.Errorf("Error parsing int style time_remaining")
	}
}
