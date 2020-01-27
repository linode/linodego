package integration

import (
	"testing"
	"time"
)

func assertDateSet(t *testing.T, compared *time.Time) {
	emptyTime := time.Time{}
	if compared == nil || *compared == emptyTime {
		t.Errorf("Expected date to be set, got %v", compared)
	}
}
