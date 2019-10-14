package timy

import "context"

// BagEngine represents a way to store persistently your events and logs
type BagEngine interface {
	RegisterNewEventRoot(c context.Context, eventRoot *EventRoot) error
	RegisterNewEventType(c context.Context, eventType *EventType) error
	RegisterNewEntry(c context.Context, entry *Entry) error
	VerifyIfEventRootExist(c context.Context, eventRootID string) error
	VerifyIfEventTypeExist(c context.Context, eventTypeID string) error
	Close() error
}
