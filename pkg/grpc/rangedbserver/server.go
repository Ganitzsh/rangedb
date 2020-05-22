package rangedbserver

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/inklabs/rangedb"
	"github.com/inklabs/rangedb/pkg/grpc/rangedbpb"
	"github.com/inklabs/rangedb/provider/inmemorystore"
)

type rangeDBServer struct {
	store rangedb.Store

	sync                sync.RWMutex
	allEventConnections map[PbRecordSender]struct{}

	broadcastMutex sync.Mutex
}

// Option defines functional option parameters for rangeDBServer.
type Option func(*rangeDBServer)

// WithStore is a functional option to inject a Store.
func WithStore(store rangedb.Store) Option {
	return func(api *rangeDBServer) {
		api.store = store
	}
}

// New constructs a new rangeDBServer.
func New(options ...Option) *rangeDBServer {
	server := &rangeDBServer{
		store:               inmemorystore.New(),
		allEventConnections: make(map[PbRecordSender]struct{}),
	}

	for _, option := range options {
		option(server)
	}

	server.initProjections()

	return server
}

func (s *rangeDBServer) initProjections() {
	s.store.Subscribe(rangedb.RecordSubscriberFunc(s.broadcastRecord))
}

func (s *rangeDBServer) Events(req *rangedbpb.EventsRequest, stream rangedbpb.RangeDB_EventsServer) error {
	for record := range s.store.EventsStartingWith(req.StartingWithEventNumber) {
		pbRecord, err := rangedbpb.ToPbRecord(record)
		if err != nil {
			return err
		}

		err = stream.Send(pbRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *rangeDBServer) EventsByStream(req *rangedbpb.EventsByStreamRequest, stream rangedbpb.RangeDB_EventsByStreamServer) error {
	for record := range s.store.EventsByStreamStartingWith(req.StreamName, req.StartingWithEventNumber) {
		pbRecord, err := rangedbpb.ToPbRecord(record)
		if err != nil {
			return err
		}

		err = stream.Send(pbRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *rangeDBServer) EventsByAggregateType(req *rangedbpb.EventsByAggregateTypeRequest, stream rangedbpb.RangeDB_EventsByAggregateTypeServer) error {
	for record := range s.store.EventsByAggregateTypesStartingWith(req.StartingWithEventNumber, req.AggregateTypes...) {
		pbRecord, err := rangedbpb.ToPbRecord(record)
		if err != nil {
			return err
		}

		err = stream.Send(pbRecord)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *rangeDBServer) SaveEvents(_ context.Context, req *rangedbpb.SaveEventsRequest) (*rangedbpb.SaveEventResponse, error) {
	eventsSaved := uint32(0)

	for _, event := range req.Events {
		var data interface{}
		err := json.Unmarshal([]byte(event.Data), &data)
		if err != nil {
			st := status.New(codes.InvalidArgument, fmt.Sprintf("unable to read event data: %v", err))
			st, _ = st.WithDetails(&rangedbpb.SaveEventFailureResponse{
				EventsSaved: eventsSaved,
			})
			return nil, st.Err()
		}

		var metadata interface{}
		if event.Metadata != "" {
			err = json.Unmarshal([]byte(event.Metadata), &metadata)
			if err != nil {
				st := status.New(codes.InvalidArgument, fmt.Sprintf("unable to read event metadata: %v", err))
				st, _ = st.WithDetails(&rangedbpb.SaveEventFailureResponse{
					EventsSaved: eventsSaved,
				})
				return nil, st.Err()
			}
		}

		err = s.store.SaveEvent(
			req.AggregateType,
			req.AggregateID,
			event.Type,
			event.ID,
			data,
			metadata,
		)
		if err != nil {
			st := status.New(codes.Internal, fmt.Sprintf("unable to save to store: %v", err))
			st, _ = st.WithDetails(&rangedbpb.SaveEventFailureResponse{
				EventsSaved: eventsSaved,
			})
			return nil, st.Err()
		}

		eventsSaved++
	}

	return &rangedbpb.SaveEventResponse{
		EventsSaved: eventsSaved,
	}, nil
}

func (s *rangeDBServer) SubscribeToEvents(req *rangedbpb.EventsRequest, stream rangedbpb.RangeDB_SubscribeToEventsServer) error {
	for record := range s.store.EventsStartingWith(req.StartingWithEventNumber) {
		pbRecord, err := rangedbpb.ToPbRecord(record)
		if err != nil {
			return err
		}

		err = stream.Send(pbRecord)
		if err != nil {
			return err
		}
	}

	s.sync.Lock()
	s.allEventConnections[stream] = struct{}{}
	s.sync.Unlock()

	<-stream.Context().Done()

	s.sync.Lock()
	delete(s.allEventConnections, stream)
	s.sync.Unlock()

	return nil
}

func (s *rangeDBServer) broadcastRecord(record *rangedb.Record) {
	s.broadcastMutex.Lock()
	go func() {
		defer s.broadcastMutex.Unlock()

		pbRecord, err := rangedbpb.ToPbRecord(record)
		if err != nil {
			//s.logger.Printf("unable to marshal record: %v", err)
			return
		}

		s.sync.RLock()
		defer s.sync.RUnlock()

		for connection := range s.allEventConnections {
			err := connection.Send(pbRecord)
			if err != nil {
				//s.logger.Printf("unable to send record to gRPC client: %v", err)
			}
		}
	}()
}

type PbRecordSender interface {
	Send(*rangedbpb.Record) error
}
