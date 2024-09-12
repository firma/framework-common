package utils

import "testing"

func TestNumberSliceContains(t *testing.T) {
	check := []int64{4, 1, 2, 3, 5}
	if NumberSliceContains(check, 14) {
		t.Error("Expected true")
	} else {
		t.Error("Expected false")
	}
}
