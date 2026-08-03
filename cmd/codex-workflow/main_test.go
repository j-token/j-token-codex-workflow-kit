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

func TestNormalizeReviewArgsKeepsStartAtValueWithFlag(t *testing.T) {
	actual := normalizeReviewArgs([]string{"prompt.md", "--start-at", "plan", "--json"})
	expected := []string{"--start-at", "plan", "--json", "prompt.md"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("normalizeReviewArgs() = %#v", actual)
	}
}

func TestReviewBaselineRoundTrip(t *testing.T) {
	documentPath := t.TempDir() + "/prompt.md"
	t.Cleanup(func() { _ = clearReviewBaseline(documentPath) })

	if err := saveReviewBaseline(documentPath, "이전 계획"); err != nil {
		t.Fatalf("saveReviewBaseline() error = %v", err)
	}
	content, err := loadReviewBaseline(documentPath)
	if err != nil || content != "이전 계획" {
		t.Fatalf("loadReviewBaseline() = %q, %v", content, err)
	}
	if err := clearReviewBaseline(documentPath); err != nil {
		t.Fatalf("clearReviewBaseline() error = %v", err)
	}
	content, err = loadReviewBaseline(documentPath)
	if err != nil || content != "" {
		t.Fatalf("cleared baseline = %q, %v", content, err)
	}
}
