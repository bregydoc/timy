package timy

// Entry define a simple entry of your event
type Entry struct {
	Type      EventType
	Modifiers map[string]interface{}
	Value     int64
}
