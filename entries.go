package timy

import "time"

// Entry define a simple entry of your event
type Entry struct {
	ID        string
	At        time.Time
	Type      EventType
	Modifiers map[string]interface{}
	Value     int64
}
