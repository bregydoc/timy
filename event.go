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
	Root        EventRoot
	Identifier  string
	CreatedAt   time.Time
	Occurrences int64
}
