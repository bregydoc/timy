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

// EventRootUpdate is a payload to update your event root
type EventRootUpdate struct {
	Name *string
}

// EventTypeUpdate is a payload to update your event root type
type EventTypeUpdate struct {
	Name        *string
	RootID      *string
	Identifier  *string
	CreatedAt   *time.Time
	Occurrences *int64
}
