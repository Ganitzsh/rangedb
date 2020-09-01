package eventstore_test

import (
	"testing"
	"time"

	"github.com/inklabs/rangedb"
	"github.com/inklabs/rangedb/pkg/clock"
	"github.com/inklabs/rangedb/provider/eventstore"
	"github.com/inklabs/rangedb/rangedbtest"
)

func Test_EventStore_VerifyStoreInterface(t *testing.T) {
	rangedbtest.VerifyStore(t, func(t *testing.T, clock clock.Clock) rangedb.Store {
		store := eventstore.New(
			"http://127.0.0.1:2113",
			"admin",
			"changeit",
			eventstore.WithClock(clock),
		)
		version := time.Now().Unix()
		store.SetVersion(version)

		t.Cleanup(func() {
			version++
			store.SetVersion(version)
		})

		return store
	})
}
