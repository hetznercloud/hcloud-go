package hcloud

import "time"

// Ptr returns a pointer to p.
func Ptr[T any](p T) *T {
	return &p
}

// String returns a pointer to the passed string s.
//
// Deprecated: Use [Ptr] instead.
func String(s string) *string { return Ptr(s) }

// Int returns a pointer to the passed integer i.
//
// Deprecated: Use [Ptr] instead.
func Int(i int) *int { return Ptr(i) }

// Bool returns a pointer to the passed bool b.
//
// Deprecated: Use [Ptr] instead.
func Bool(b bool) *bool { return Ptr(b) }

// Duration returns a pointer to the passed time.Duration d.
//
// Deprecated: Use [Ptr] instead.
func Duration(d time.Duration) *time.Duration { return Ptr(d) }

// copySlice returns a copy of the slice s. If s is nil, copySlice returns nil.
func copySlice[T any](s []T) []T {
	if s == nil {
		return nil
	}
	clone := make([]T, len(s))
	copy(clone, s)
	return clone
}
