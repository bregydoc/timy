package timy

import (
	"context"

	badger "github.com/dgraph-io/badger"
	"github.com/influxdata/influxdb-client-go"
)

// InfluxBag is an influx bag implementation
type InfluxBag struct {
	dict   *badger.DB
	client *influxdb.Client
}

// InfluxBagConfig is a wrapper for your bag configurations
// TODO: Improve the config of badger dictionary, accept more options
type InfluxBagConfig struct {
	InfluxConnection string
	InfluxToken      string
	BadgerFilepath   string
}

// DefaultInfluxBagConfig creates a new default configuration for your influx bag
func DefaultInfluxBagConfig(influxConn, influxToken string) *InfluxBagConfig {
	return &InfluxBagConfig{
		InfluxConnection: influxConn,
		InfluxToken:      influxToken,
		BadgerFilepath:   "./dict.db",
	}
}

// NewInfluxBag generates a new bag instance
func NewInfluxBag(config *InfluxBagConfig) (*InfluxBag, error) {
	dict, err := badger.Open(badger.DefaultOptions(config.BadgerFilepath))
	if err != nil {
		return nil, err
	}

	client, err := influxdb.New(config.InfluxConnection, config.InfluxToken)
	if err != nil {
		return nil, err
	}

	return &InfluxBag{
		dict:   dict,
		client: client,
	}, nil
}

// RegisterNewEventRoot implements the bag interface
func (bag *InfluxBag) RegisterNewEventRoot(c context.Context, eventRoot *EventRoot) error {
	id := generateNewID()
	eventRoot.ID = id
	return bag.indexEventRoot(id, eventRoot.Name)
}

// RegisterNewEventType implements the bag interface
func (bag *InfluxBag) RegisterNewEventType(c context.Context, eventType *EventType) error {
	id := generateNewID()
	eventType.ID = id

	return bag.indexEventType(
		id, eventType.Name,
		eventType.RootID,
		eventType.Identifier,
		eventType.CreatedAt,
		eventType.Occurrences,
	)
}

// UpdateEventRoot implements the bag interface
func (bag *InfluxBag) UpdateEventRoot(c context.Context, eventRootID string, update *EventRootUpdate) error {
	eventRoot, err := bag.getEventRoot(eventRootID)
	if err != nil {
		return err
	}

	if update == nil {
		return nil
	}

	if update.Name != nil {
		eventRoot.Name = *update.Name
	}

	return bag.indexEventRoot(eventRootID, eventRoot.Name)

}

// UpdateEventType implements the bag interface
func (bag *InfluxBag) UpdateEventType(c context.Context, eventTypeID string, update *EventTypeUpdate) error {
	eventType, err := bag.getEventType(eventTypeID)
	if err != nil {
		return err
	}

	if update == nil {
		return nil
	}

	if update.Name != nil {
		eventType.Name = *update.Name
	}

	if update.RootID != nil {
		eventType.RootID = *update.RootID
	}

	if update.Identifier != nil {
		eventType.Identifier = *update.Identifier
	}

	if update.CreatedAt != nil {
		eventType.CreatedAt = *update.CreatedAt
	}

	if update.Occurrences != nil {
		eventType.Occurrences = *update.Occurrences
	}

	return bag.indexEventType(
		eventTypeID,
		eventType.Name,
		eventType.RootID,
		eventType.Identifier,
		eventType.CreatedAt,
		eventType.Occurrences,
	)
}

// RegisterNewEntry implements the bag interface
func (bag *InfluxBag) RegisterNewEntry(c context.Context, entry *Entry) error {
	return bag.newInfluxMetricRow(c, bag.client, entry, false)
}

// VerifyIfEventRootExist implements the bag interface
func (bag *InfluxBag) VerifyIfEventRootExist(c context.Context, eventRootID string) error {
	_, err := bag.getEventRoot(eventRootID)
	return err
}

// VerifyIfEventTypeExist implements the bag interface
func (bag *InfluxBag) VerifyIfEventTypeExist(c context.Context, eventTypeID string) error {
	_, err := bag.getEventType(eventTypeID)
	return err
}

// Close implements the bag interface
func (bag *InfluxBag) Close() error {
	if err := bag.Close(); err != nil {
		return err
	}

	return bag.client.Close()
}
