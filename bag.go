package timy

import "context"

// BagEngine represents a way to store persistently your events and logs
type BagEngine interface {
	RegisterNewEventRoot(c context.Context, eventRoot *EventRoot) error
	UpdateEventRoot(c context.Context, eventRootID string, update *EventRootUpdate) error
	RegisterNewEventType(c context.Context, eventType *EventType) error
	UpdateEventType(c context.Context, eventTypeID string, update *EventTypeUpdate) error
	RegisterNewEntry(c context.Context, entry *Entry) error
	VerifyIfEventRootExist(c context.Context, eventRootID string) error
	VerifyIfEventTypeExist(c context.Context, eventTypeID string) error
	Close() error
}
