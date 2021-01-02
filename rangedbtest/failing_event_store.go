package rangedbtest

import (
	"context"
	"fmt"

	"github.com/inklabs/rangedb"
)

type failingEventStore struct{}

// NewFailingEventStore constructs a failing event store for testing.
func NewFailingEventStore() *failingEventStore {
	return &failingEventStore{}
}

func (f failingEventStore) Bind(_ ...rangedb.Event) {}

func (f failingEventStore) EventsStartingWith(_ context.Context, _ uint64) rangedb.RecordIterator {
	return getClosedIterator()
}

func (f failingEventStore) EventsByAggregateTypesStartingWith(_ context.Context, _ uint64, _ ...string) rangedb.RecordIterator {
	return getClosedIterator()
}

func (f failingEventStore) EventsByStreamStartingWith(_ context.Context, _ uint64, _ string) rangedb.RecordIterator {
	return getClosedIterator()
}

func (f failingEventStore) OptimisticSave(_ context.Context, _ uint64, _ ...*rangedb.EventRecord) error {
	return fmt.Errorf("failingEventStore.OptimisticSave")
}

func (f failingEventStore) Save(_ context.Context, _ ...*rangedb.EventRecord) error {
	return fmt.Errorf("failingEventStore.Save")
}

func (f failingEventStore) Subscribe(_ context.Context, _ ...rangedb.RecordSubscriber) error {
	return nil
}

func (f failingEventStore) SubscribeStartingWith(_ context.Context, _ uint64, _ ...rangedb.RecordSubscriber) error {
	return nil
}

func (f failingEventStore) TotalEventsInStream(_ context.Context, _ string) (uint64, error) {
	return 0, fmt.Errorf("failingEventStore.TotalEventsInStream")
}

func getClosedIterator() rangedb.RecordIterator {
	recordResults := make(chan rangedb.ResultRecord)
	close(recordResults)
	return rangedb.NewRecordIterator(recordResults)
}
