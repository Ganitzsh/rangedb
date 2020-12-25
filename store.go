package rangedb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const Version = "0.5.1-dev"

// Record contains event data and metadata.
type Record struct {
	AggregateType        string      `msgpack:"a" json:"aggregateType"`
	AggregateID          string      `msgpack:"i" json:"aggregateID"`
	GlobalSequenceNumber uint64      `msgpack:"g" json:"globalSequenceNumber"`
	StreamSequenceNumber uint64      `msgpack:"s" json:"sequenceNumber"`
	InsertTimestamp      uint64      `msgpack:"u" json:"insertTimestamp"`
	EventID              string      `msgpack:"e" json:"eventID"`
	EventType            string      `msgpack:"t" json:"eventType"`
	Data                 interface{} `msgpack:"d" json:"data"`
	Metadata             interface{} `msgpack:"m" json:"metadata"`
}

type EventRecord struct {
	Event    Event
	Metadata interface{}
}

type EventBinder interface {
	Bind(events ...Event)
}

// Store is the interface that stores and retrieves event records.
type Store interface {
	EventBinder
	EventsStartingWith(ctx context.Context, globalSequenceNumber uint64) <-chan *Record
	EventsByAggregateTypesStartingWith(ctx context.Context, globalSequenceNumber uint64, aggregateTypes ...string) <-chan *Record
	EventsByStreamStartingWith(ctx context.Context, streamSequenceNumber uint64, streamName string) <-chan *Record
	OptimisticSave(expectedStreamSequenceNumber uint64, eventRecords ...*EventRecord) error
	Save(eventRecords ...*EventRecord) error
	Subscribe(subscribers ...RecordSubscriber)
	SubscribeStartingWith(ctx context.Context, globalSequenceNumber uint64, subscribers ...RecordSubscriber)
	TotalEventsInStream(streamName string) uint64
}

// Event is the interface that defines the required event methods.
type Event interface {
	AggregateMessage
	EventType() string
}

// AggregateMessage is the interface that supports building an event stream name.
type AggregateMessage interface {
	AggregateID() string
	AggregateType() string
}

// The RecordSubscriberFunc type is an adapter to allow the use of
// ordinary functions as record subscribers. If f is a function
// with the appropriate signature, RecordSubscriberFunc(f) is a
// Handler that calls f.
type RecordSubscriberFunc func(*Record)

func (f RecordSubscriberFunc) Accept(record *Record) {
	f(record)
}

// RecordSubscriber is the interface that defines how a projection receives Records.
type RecordSubscriber interface {
	Accept(record *Record)
}

// GetEventStream returns the stream name for an event.
func GetEventStream(message AggregateMessage) string {
	return GetStream(message.AggregateType(), message.AggregateID())
}

// GetStream returns the stream name for an aggregateType and aggregateID.
func GetStream(aggregateType, aggregateID string) string {
	return fmt.Sprintf("%s!%s", aggregateType, aggregateID)
}

// ParseStream returns the aggregateType and aggregateID for a stream name.
func ParseStream(streamName string) (aggregateType, aggregateID string) {
	pieces := strings.Split(streamName, "!")
	return pieces[0], pieces[1]
}

// ReplayEvents applies all events to each subscriber.
func ReplayEvents(store Store, globalSequenceNumber uint64, subscribers ...RecordSubscriber) {
	for record := range store.EventsStartingWith(context.Background(), globalSequenceNumber) {
		for _, subscriber := range subscribers {
			subscriber.Accept(record)
		}
	}
}

// ReadNRecords reads up to N records from the channel returned by f into a slice
func ReadNRecords(totalEvents uint64, f func(context.Context) <-chan *Record) []*Record {
	var records []*Record
	ctx, done := context.WithCancel(context.Background())
	cnt := uint64(0)
	for record := range f(ctx) {
		cnt++
		if cnt > totalEvents {
			break
		}

		records = append(records, record)
	}
	done()

	return records
}

type rawEvent struct {
	aggregateType string
	aggregateID   string
	eventType     string
	data          interface{}
}

func NewRawEvent(aggregateType, aggregateID, eventType string, data interface{}) Event {
	return &rawEvent{
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
		eventType:     eventType,
		data:          data,
	}
}

func (e rawEvent) AggregateID() string {
	return e.aggregateID
}

func (e rawEvent) AggregateType() string {
	return e.aggregateType
}

func (e rawEvent) EventType() string {
	return e.eventType
}

func (e rawEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.data)
}
