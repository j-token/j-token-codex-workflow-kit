package main

import (
	"reflect"
	"testing"
)

func TestNormalizeReviewArgsAllowsFlagsAfterPath(t *testing.T) {
	actual := normalizeReviewArgs([]string{"prompt.md", "--json", "--no-open"})
	expected := []string{"--json", "--no-open", "prompt.md"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("normalizeReviewArgs() = %#v", actual)
	}
}
