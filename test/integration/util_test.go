package integration

import (
	"strconv"
	"testing"
	"time"
)

func assertDateSet(t *testing.T, compared *time.Time) {
	emptyTime := time.Time{}
	if compared == nil || *compared == emptyTime {
		t.Errorf("Expected date to be set, got %v", compared)
	}
}

func assertSliceContains[T comparable](t *testing.T, slice []T, target T) {
	for _, v := range slice {
		if v == target {
			return
		}
	}

	t.Fatalf("value %v not found in slice", target)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// return the current nanosecond in string type as a unique text.
func getUniqueText() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

// getIntersection returns common items of two data structures (e.g. slices)
func getIntersection[T comparable](a, b []T) []T {
	set := make(map[T]struct{})
	var result []T

	for _, item := range a {
		set[item] = struct{}{}
	}

	for _, item := range b {
		if _, ok := set[item]; ok {
			result = append(result, item)
		}
	}
	return result
}
