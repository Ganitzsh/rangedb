package eventstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/inklabs/rangedb"
	"github.com/inklabs/rangedb/pkg/clock"
	"github.com/inklabs/rangedb/pkg/clock/provider/systemclock"
	"github.com/inklabs/rangedb/pkg/shortuuid"
	"github.com/inklabs/rangedb/provider/jsonrecordserializer"
)

type eventStore struct {
	clock      clock.Clock
	serializer rangedb.RecordSerializer
	httpClient *http.Client
	baseURI    string
	username   string
	password   string
	version    int64
}

// Option defines functional option parameters for eventStore.
type Option func(*eventStore)

func WithClock(clock clock.Clock) Option {
	return func(store *eventStore) {
		store.clock = clock
	}
}

// WithSerializer is a functional option to inject a RecordSerializer.
func WithSerializer(serializer rangedb.RecordSerializer) Option {
	return func(store *eventStore) {
		store.serializer = serializer
	}
}

// New constructs an eventStore.
func New(baseURI, username, password string, options ...Option) *eventStore {
	s := &eventStore{
		clock:      systemclock.New(),
		serializer: jsonrecordserializer.New(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURI:  baseURI,
		username: username,
		password: password,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *eventStore) Bind(events ...rangedb.Event) {
	s.serializer.Bind(events...)
}

func (s *eventStore) EventsStartingWith(ctx context.Context, eventNumber uint64) <-chan *rangedb.Record {
	return s.EventsByStreamStartingWith(ctx, eventNumber, "$all")
}

func (s *eventStore) EventsByAggregateTypesStartingWith(ctx context.Context, eventNumber uint64, aggregateTypes ...string) <-chan *rangedb.Record {
	records := make(chan *rangedb.Record)

	go func() {
		defer close(records)

		aggregateTypesMap := make(map[string]struct{}, len(aggregateTypes))
		for _, aggregateType := range aggregateTypes {
			aggregateTypesMap[aggregateType] = struct{}{}
		}

		allRecords := s.EventsByStreamStartingWith(ctx, eventNumber, "$all")
		for record := range allRecords {
			if _, ok := aggregateTypesMap[record.AggregateType]; ok {
				records <- record
			}
		}
	}()

	return records
}

func (s *eventStore) EventsByStreamStartingWith(ctx context.Context, eventNumber uint64, streamName string) <-chan *rangedb.Record {
	records := make(chan *rangedb.Record)

	go func() {
		defer close(records)

		firstPage := fmt.Sprintf("%s/streams/%s/0/forward/20", s.baseURI, s.streamName(streamName))

		count := uint64(0)
		uri := firstPage
		for {
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, uri+"?embed=body", nil)
			if err != nil {
				log.Printf("unable to create request: %v", err)
				return
			}

			request.SetBasicAuth(s.username, s.password)
			request.Header.Add("Accept", "application/vnd.eventstore.atom+json")

			response, err := s.httpClient.Do(request)
			if err != nil {
				log.Printf("unable to send request: %v", err)
				return
			}

			if response.StatusCode != http.StatusOK {
				log.Printf("error\nrequest: %#v\nresponse: %#v\n\n", request.URL.String(), response)
				return
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf("error reading response body: %#v", err)
				return
			}
			_ = response.Body.Close()
			var streamResponse StreamResponse
			err = json.Unmarshal(body, &streamResponse)
			if err != nil {
				log.Printf("error unmarshaling response body: %v", err)
				return
			}

			for i := len(streamResponse.Entries) - 1; i >= 0; i-- {
				if count >= eventNumber {
					entry := streamResponse.Entries[i]

					record, err := s.serializer.Deserialize([]byte(entry.Data))
					if err != nil {
						log.Printf("error decoding event record: %v", err)
						return
					}

					record.GlobalSequenceNumber = uint64(entry.EventNumber)
					record.StreamSequenceNumber = uint64(entry.PositionEventNumber)

					records <- record
				}
				count++
			}

			uri, err = streamResponse.NextLink()
			if err != nil {
				log.Printf("no more pages: %#v", streamResponse.Links)
				return
			}
		}
	}()

	return records
}

func (s *eventStore) Save(event rangedb.Event, metadata interface{}) error {
	return s.SaveEvent(
		event.AggregateType(),
		event.AggregateID(),
		event.EventType(),
		shortuuid.New().String(),
		event,
		metadata,
	)
}

func (s *eventStore) SaveEvent(aggregateType, aggregateID, eventType, eventID string, event, metadata interface{}) error {
	url := fmt.Sprintf("%s/streams/%s", s.baseURI, s.streamName(rangedb.GetStream(aggregateType, aggregateID)))

	record := &rangedb.Record{
		AggregateType:        aggregateType,
		AggregateID:          aggregateID,
		GlobalSequenceNumber: 0,
		StreamSequenceNumber: 0,
		EventType:            eventType,
		EventID:              eventID,
		InsertTimestamp:      uint64(s.clock.Now().Unix()),
		Data:                 event,
		Metadata:             metadata,
	}

	recordData, err := s.serializer.Serialize(record)
	if err != nil {
		return err
	}

	var eventMetadata []byte

	if metadata != nil {
		eventMetadata, err = json.Marshal(metadata)
		if err != nil {
			return err
		}
	}

	newEvents := []NewEvent{
		{
			EventID:   eventID,
			EventType: eventType,
			Data:      string(recordData),
			Metadata:  string(eventMetadata),
		},
	}

	requestBody := &bytes.Buffer{}
	err = json.NewEncoder(requestBody).Encode(newEvents)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return err
	}

	request.SetBasicAuth(s.username, s.password)
	request.Header.Add("Content-Type", "application/vnd.eventstore.events+json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		log.Printf("error response: %#v", response)
		return fmt.Errorf("unable to save event")
	}

	return nil
}

func (s *eventStore) Subscribe(subscribers ...rangedb.RecordSubscriber) {

}

func (s *eventStore) SubscribeStartingWith(ctx context.Context, eventNumber uint64, subscribers ...rangedb.RecordSubscriber) {

}

func (s *eventStore) TotalEventsInStream(streamName string) uint64 {
	panic("implement me")
}

func (s *eventStore) streamName(name string) string {
	return fmt.Sprintf("%d-%s", s.version, name)
}

func (s *eventStore) SetVersion(version int64) {
	s.version = version
}
