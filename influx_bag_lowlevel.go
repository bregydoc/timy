package timy

import (
	"encoding/binary"
	"time"
)

const eventRootIndexKey = "eventroot."
const eventTypeIndexKey = "eventtype."

func (bag *InfluxBag) eventRootKey(eventRootID string) []byte {
	return []byte(eventRootIndexKey + eventRootID)
}

func (bag *InfluxBag) eventTypeNameKey(eventTypeID string) []byte {
	return []byte(eventTypeIndexKey + "name." + eventTypeID)
}

func (bag *InfluxBag) eventTypeRootIDKey(eventTypeID string) []byte {
	return []byte(eventTypeIndexKey + "root." + eventTypeID)
}

func (bag *InfluxBag) eventTypeIdentifierKey(eventTypeID string) []byte {
	return []byte(eventTypeIndexKey + "ident." + eventTypeID)
}

func (bag *InfluxBag) eventTypeCreatedAtKey(eventTypeID string) []byte {
	return []byte(eventTypeIndexKey + "createdat." + eventTypeID)
}

func (bag *InfluxBag) eventTypeOccurrencesAtKey(eventTypeID string) []byte {
	return []byte(eventTypeIndexKey + "occurrences." + eventTypeID)
}

func (bag *InfluxBag) indexEventRoot(eventRootID string, eventRootName string) error {
	txn := bag.dict.NewTransaction(true)
	defer txn.Discard()

	if err := txn.Set(bag.eventRootKey(eventRootID), []byte(eventRootName)); err != nil {
		return err
	}

	return txn.Commit()
}

func (bag *InfluxBag) getEventRoot(eventRootID string) (*EventRoot, error) {
	txn := bag.dict.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(bag.eventRootKey(eventRootID))
	if err != nil {
		return nil, err
	}

	var name string
	err = item.Value(func(val []byte) error {
		name = string(val)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &EventRoot{
		ID:   eventRootID,
		Name: name,
	}, nil
}

func (bag *InfluxBag) indexEventType(eventTypeID string, name, rootID, identifier string, createdAt time.Time, occurrences int64) error {
	txn := bag.dict.NewTransaction(true)
	defer txn.Discard()

	createdAtString := createdAt.Format(time.RFC3339)

	binOccur := make([]byte, 8)
	binary.LittleEndian.PutUint64(binOccur, uint64(occurrences)) // So good :D

	if err := txn.Set(bag.eventTypeNameKey(eventTypeID), []byte(name)); err != nil {
		return err
	}

	if err := txn.Set(bag.eventTypeRootIDKey(eventTypeID), []byte(rootID)); err != nil {
		return err
	}

	if err := txn.Set(bag.eventTypeIdentifierKey(eventTypeID), []byte(identifier)); err != nil {
		return err
	}

	if err := txn.Set(bag.eventTypeCreatedAtKey(eventTypeID), []byte(createdAtString)); err != nil {
		return err
	}

	if err := txn.Set(bag.eventTypeOccurrencesAtKey(eventTypeID), binOccur); err != nil {
		return err
	}

	return txn.Commit()
}

func (bag *InfluxBag) getEventType(eventTypeID string) (*EventType, error) {
	txn := bag.dict.NewTransaction(false)
	defer txn.Discard()

	nameItem, err := txn.Get(bag.eventTypeNameKey(eventTypeID))
	if err != nil {
		return nil, err
	}

	rootIDItem, err := txn.Get(bag.eventTypeRootIDKey(eventTypeID))
	if err != nil {
		return nil, err
	}

	identifierItem, err := txn.Get(bag.eventTypeIdentifierKey(eventTypeID))
	if err != nil {
		return nil, err
	}

	createAtItem, err := txn.Get(bag.eventTypeCreatedAtKey(eventTypeID))
	if err != nil {
		return nil, err
	}

	occurItem, err := txn.Get(bag.eventTypeOccurrencesAtKey(eventTypeID))
	if err != nil {
		return nil, err
	}

	var name string
	var rootID string
	var identifier string
	var createdAt time.Time
	var occurrences int64

	if err = nameItem.Value(func(val []byte) error {
		name = string(val)
		return nil
	}); err != nil {
		return nil, err
	}

	if err = rootIDItem.Value(func(val []byte) error {
		rootID = string(val)
		return nil
	}); err != nil {
		return nil, err
	}

	if err = identifierItem.Value(func(val []byte) error {
		identifier = string(val)
		return nil
	}); err != nil {
		return nil, err
	}

	if err = createAtItem.Value(func(val []byte) error {
		createdAt, err = time.Parse(time.RFC3339, string(val))
		return err
	}); err != nil {
		return nil, err
	}

	if err = occurItem.Value(func(val []byte) error {
		occurrences = int64(binary.LittleEndian.Uint64(val))
		return nil
	}); err != nil {
		return nil, err
	}

	return &EventType{
		ID:          eventTypeID,
		Name:        name,
		RootID:      rootID,
		Identifier:  identifier,
		CreatedAt:   createdAt,
		Occurrences: occurrences,
	}, nil
}
