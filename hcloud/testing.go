package hcloud

import (
	"testing"
	"time"
)

const apiTimestampFormat = time.RFC3339

func mustParseTime(t *testing.T, layout, value string) time.Time {
	t.Helper()

	ts, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("parse time: layout %v: value %v: %v", layout, value, err)
	}
	return ts
}
