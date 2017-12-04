package api

import "time"

// Action represents an action in the Hetzner Cloud.
type Action struct {
	ID           int
	Status       string
	Command      string
	Progress     int
	Started      time.Time
	Finished     time.Time
	ErrorCode    string
	ErrorMessage string
}
