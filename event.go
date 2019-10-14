package timy

import "time"

// EventRoot describes a simple event root name
type EventRoot struct {
	ID   string
	Name string
}

// EventType represents a type of event, it has a parent
type EventType struct {
	ID          string
	Name        string
	RootID      string
	Identifier  string
	CreatedAt   time.Time
	Occurrences int64
}
