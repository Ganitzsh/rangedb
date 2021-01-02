package rangedbtest

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/rangedb"
)

// StoreBenchmark benchmarks the rangedb.Store interface.
func StoreBenchmark(b *testing.B, newStore func(b *testing.B) rangedb.Store) {
	b.Helper()

	for _, totalEvents := range []int{1, 5, 10, 50} {
		eventRecords := getNEvents(totalEvents, "eb4b1c61fa344272a61e039cc4247258")
		b.Run(fmt.Sprintf("Save %d at a time", totalEvents), func(b *testing.B) {
			store := newStore(b)
			ctx := TimeoutContext(b)
			for i := 0; i < b.N; i++ {
				err := store.Save(ctx, eventRecords...)
				if err != nil {
					require.NoError(b, err)
				}
			}
		})

		b.Run(fmt.Sprintf("Optimistic Save %d at a time", totalEvents), func(b *testing.B) {
			store := newStore(b)
			ctx := TimeoutContext(b)
			for i := 0; i < b.N; i++ {
				err := store.OptimisticSave(ctx, uint64(i*totalEvents), eventRecords...)
				if err != nil {
					require.NoError(b, err)
				}
			}
		})
	}

	b.Run("Queries", func(b *testing.B) {
		store := newStore(b)
		const (
			totalEvents                  = 10000
			eventsPerStream              = 10
			totalEventsToRead            = 1000
			startingGlobalSequenceNumber = totalEvents - totalEventsToRead
		)

		aggregateID := ""
		ctx := TimeoutContext(b)
		for i := 0; i < (totalEvents / eventsPerStream); i++ {
			aggregateID = uuid.New().String()
			eventRecords := getNEvents(eventsPerStream, aggregateID)
			err := store.Save(ctx, eventRecords...)
			require.NoError(b, err)
		}

		b.Run("EventsStartingWith", func(b *testing.B) {
			ctx := TimeoutContext(b)
			for i := 0; i < b.N; i++ {
				recordIterator := store.EventsStartingWith(ctx, startingGlobalSequenceNumber)
				cnt := 0
				for recordIterator.Next() {
					if recordIterator.Err() != nil {
						break
					}
					cnt++
				}
				require.Equal(b, totalEventsToRead, cnt)
			}
		})

		b.Run("EventsByAggregateTypesStartingWith", func(b *testing.B) {
			ctx := TimeoutContext(b)
			for i := 0; i < b.N; i++ {
				recordIterator := store.EventsByAggregateTypesStartingWith(ctx, startingGlobalSequenceNumber, ThingWasDone{}.AggregateType())
				cnt := 0
				for recordIterator.Next() {
					if recordIterator.Err() != nil {
						break
					}
					cnt++
				}
				require.Equal(b, totalEventsToRead, cnt)
			}
		})

		b.Run("EventsByStreamStartingWith", func(b *testing.B) {
			ctx := TimeoutContext(b)
			stream := rangedb.GetStream(ThingWasDone{}.AggregateType(), aggregateID)
			for i := 0; i < b.N; i++ {
				recordIterator := store.EventsByStreamStartingWith(ctx, 0, stream)
				cnt := 0
				for recordIterator.Next() {
					if recordIterator.Err() != nil {
						break
					}
					cnt++
				}
				require.Equal(b, eventsPerStream, cnt)
			}
		})

		b.Run("TotalEventsInStream", func(b *testing.B) {
			stream := rangedb.GetStream(ThingWasDone{}.AggregateType(), aggregateID)
			for i := 0; i < b.N; i++ {
				totalEventsInStream := store.TotalEventsInStream(stream)
				require.Equal(b, uint64(eventsPerStream), totalEventsInStream)
			}
		})
	})
}

func getNEvents(n int, aggregateID string) []*rangedb.EventRecord {
	eventRecords := make([]*rangedb.EventRecord, n)

	for i := 0; i < n; i++ {
		eventRecords[i] = &rangedb.EventRecord{
			Event: &ThingWasDone{
				ID:     aggregateID,
				Number: i,
			},
			Metadata: nil,
		}
	}

	return eventRecords
}
