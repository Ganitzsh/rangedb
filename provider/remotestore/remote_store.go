package remotestore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"google.golang.org/grpc"

	"github.com/inklabs/rangedb"
	"github.com/inklabs/rangedb/pkg/grpc/rangedbpb"
	"github.com/inklabs/rangedb/pkg/rangedberror"
	"github.com/inklabs/rangedb/provider/jsonrecordserializer"
)

// JsonSerializer defines the interface to bind events and identify event types.
type JsonSerializer interface {
	rangedb.EventBinder
	rangedb.EventTypeIdentifier
}

// PbRecordReceiver defines the interface to receive a protobuf record.
type PbRecordReceiver interface {
	Recv() (*rangedbpb.Record, error)
}

type remoteStore struct {
	serializer JsonSerializer
	client     rangedbpb.RangeDBClient

	subscriberMux sync.RWMutex
	subscribers   []rangedb.RecordSubscriber
}

// New constructs a new rangedb.Store client that communicates with a remote gRPC backend.
func New(conn *grpc.ClientConn) *remoteStore {
	client := rangedbpb.NewRangeDBClient(conn)
	return &remoteStore{
		serializer: jsonrecordserializer.New(),
		client:     client,
	}
}

func (s *remoteStore) Bind(events ...rangedb.Event) {
	s.serializer.Bind(events...)
}

func (s *remoteStore) EventsStartingWith(ctx context.Context, globalSequenceNumber uint64) <-chan *rangedb.Record {
	request := &rangedbpb.EventsRequest{
		GlobalSequenceNumber: globalSequenceNumber,
	}

	events, err := s.client.Events(ctx, request)
	if err != nil {
		return closedChannel()
	}

	return s.readRecords(ctx, events)
}

func (s *remoteStore) EventsByAggregateTypesStartingWith(ctx context.Context, globalSequenceNumber uint64, aggregateTypes ...string) <-chan *rangedb.Record {
	request := &rangedbpb.EventsByAggregateTypeRequest{
		AggregateTypes:       aggregateTypes,
		GlobalSequenceNumber: globalSequenceNumber,
	}

	events, err := s.client.EventsByAggregateType(ctx, request)
	if err != nil {
		return closedChannel()
	}

	return s.readRecords(ctx, events)

}

func (s *remoteStore) EventsByStreamStartingWith(ctx context.Context, streamSequenceNumber uint64, streamName string) <-chan *rangedb.Record {
	request := &rangedbpb.EventsByStreamRequest{
		StreamName:           streamName,
		StreamSequenceNumber: streamSequenceNumber,
	}

	events, err := s.client.EventsByStream(ctx, request)
	if err != nil {
		return closedChannel()
	}

	return s.readRecords(ctx, events)
}

func (s *remoteStore) OptimisticSave(expectedStreamSequenceNumber uint64, eventRecords ...*rangedb.EventRecord) error {
	var aggregateType, aggregateID string

	var events []*rangedbpb.Event
	for _, eventRecord := range eventRecords {
		if aggregateType != "" && aggregateType != eventRecord.Event.AggregateType() {
			return fmt.Errorf("unmatched aggregate type")
		}

		if aggregateID != "" && aggregateID != eventRecord.Event.AggregateID() {
			return fmt.Errorf("unmatched aggregate ID")
		}

		aggregateType = eventRecord.Event.AggregateType()
		aggregateID = eventRecord.Event.AggregateID()

		jsonData, err := json.Marshal(eventRecord.Event)
		if err != nil {
			return err
		}

		jsonMetadata, err := json.Marshal(eventRecord.Metadata)
		if err != nil {
			return err
		}

		events = append(events, &rangedbpb.Event{
			Type:     eventRecord.Event.EventType(),
			Data:     string(jsonData),
			Metadata: string(jsonMetadata),
		})
	}

	request := &rangedbpb.OptimisticSaveRequest{
		AggregateType:                aggregateType,
		AggregateID:                  aggregateID,
		Events:                       events,
		ExpectedStreamSequenceNumber: expectedStreamSequenceNumber,
	}

	_, err := s.client.OptimisticSave(context.Background(), request)
	if err != nil {
		if strings.Contains(err.Error(), "unable to save to store: unexpected sequence number") {
			return rangedberror.NewUnexpectedSequenceNumberFromString(err.Error())
		}

		return err
	}

	return nil
}

func (s *remoteStore) Save(eventRecords ...*rangedb.EventRecord) error {
	var aggregateType, aggregateID string

	var events []*rangedbpb.Event
	for _, eventRecord := range eventRecords {
		if aggregateType != "" && aggregateType != eventRecord.Event.AggregateType() {
			return fmt.Errorf("unmatched aggregate type")
		}

		if aggregateID != "" && aggregateID != eventRecord.Event.AggregateID() {
			return fmt.Errorf("unmatched aggregate ID")
		}

		aggregateType = eventRecord.Event.AggregateType()
		aggregateID = eventRecord.Event.AggregateID()

		jsonData, err := json.Marshal(eventRecord.Event)
		if err != nil {
			return err
		}

		jsonMetadata, err := json.Marshal(eventRecord.Metadata)
		if err != nil {
			return err
		}

		events = append(events, &rangedbpb.Event{
			Type:     eventRecord.Event.EventType(),
			Data:     string(jsonData),
			Metadata: string(jsonMetadata),
		})
	}

	request := &rangedbpb.SaveRequest{
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Events:        events,
	}

	_, err := s.client.Save(context.Background(), request)
	if err != nil {
		if strings.Contains(err.Error(), "unable to save to store: unexpected sequence number") {
			return rangedberror.NewUnexpectedSequenceNumberFromString(err.Error())
		}

		return err
	}

	return nil
}

func (s *remoteStore) SubscribeStartingWith(ctx context.Context, globalSequenceNumber uint64, subscribers ...rangedb.RecordSubscriber) {
	rangedb.ReplayEvents(s, globalSequenceNumber, subscribers...)

	select {
	case <-ctx.Done():
		return
	default:
		s.Subscribe(subscribers...)
	}
}

func (s *remoteStore) TotalEventsInStream(streamName string) uint64 {
	request := &rangedbpb.TotalEventsInStreamRequest{
		StreamName: streamName,
	}

	response, err := s.client.TotalEventsInStream(context.Background(), request)
	if err != nil {
		return 0
	}

	return response.TotalEvents
}

func (s *remoteStore) readRecords(ctx context.Context, events PbRecordReceiver) <-chan *rangedb.Record {
	records := make(chan *rangedb.Record)

	go func() {
		for {
			pbRecord, err := events.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("failed to get record: %v", err)
				break
			}

			record, err := rangedbpb.ToRecord(pbRecord, s.serializer)
			if err != nil {
				log.Printf("failed converting to record: %v", err)
				continue
			}

			select {
			case <-ctx.Done():
				break
			case records <- record:
			}
		}

		close(records)
	}()

	return records
}

func (s *remoteStore) Subscribe(subscribers ...rangedb.RecordSubscriber) {
	s.subscriberMux.Lock()
	if len(s.subscribers) == 0 {
		s.listenForEvents()
	}
	s.subscribers = append(s.subscribers, subscribers...)
	s.subscriberMux.Unlock()
}

func (s *remoteStore) listenForEvents() {
	request := &rangedbpb.SubscribeToLiveEventsRequest{}

	ctx := context.Background()
	events, err := s.client.SubscribeToLiveEvents(ctx, request)
	if err != nil {
		log.Printf("failed to subscribe: %v", err)
		return
	}

	go func() {
		for record := range s.readRecords(ctx, events) {
			s.subscriberMux.RLock()
			for _, subscriber := range s.subscribers {
				subscriber.Accept(record)
			}
			s.subscriberMux.RUnlock()
		}
	}()
}

func closedChannel() <-chan *rangedb.Record {
	records := make(chan *rangedb.Record)
	close(records)
	return records
}
