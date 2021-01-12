package chat_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/inklabs/rangedb/examples/chat"
	"github.com/inklabs/rangedb/provider/inmemorystore"
	"github.com/inklabs/rangedb/rangedbtest"
	"github.com/inklabs/rangedb/rangedbtest/bdd"
)

const (
	roomID      = "e90862685ed44bb7b656bd0b1d8f6cba"
	userID      = "c5b42a2e809e4703b5cad68de92710df"
	roomName    = "general"
	message     = "Hello, World!"
	warnMessage = "you have been warned"
	badMessage  = "Dagnabit!"
)

func Test_OnBoardUser(t *testing.T) {
	t.Run("on-boards user", newTestCase().
		Given().
		When(chat.OnBoardUser{
			UserID: userID,
			Name:   "John",
		}).
		Then(chat.UserWasOnBoarded{
			UserID: userID,
			Name:   "John",
		}))

	t.Run("fails due to existing user", newTestCase().
		Given(chat.UserWasOnBoarded{
			UserID: userID,
			Name:   "John",
		}).
		When(chat.OnBoardUser{
			UserID: userID,
			Name:   "Jane",
		}).
		Then())
}

func Test_StartRoom(t *testing.T) {
	t.Run("on-boards room", newTestCase().
		Given().
		When(chat.OnBoardRoom{
			RoomID:   roomID,
			UserID:   userID,
			RoomName: roomName,
		}).
		Then(chat.RoomWasOnBoarded{
			RoomID:   roomID,
			UserID:   userID,
			RoomName: roomName,
		}))
}

func Test_JoinRoom(t *testing.T) {
	t.Run("joins room", newTestCase().
		Given(chat.RoomWasOnBoarded{
			RoomID:   roomID,
			UserID:   userID,
			RoomName: roomName,
		}).
		When(chat.JoinRoom{
			RoomID: roomID,
			UserID: userID,
		}).
		Then(chat.RoomWasJoined{
			RoomID: roomID,
			UserID: userID,
		}))

	t.Run("fails to join invalid room", newTestCase().
		Given().
		When(chat.JoinRoom{
			RoomID: roomID,
			UserID: userID,
		}).
		Then())

	t.Run("fails to join room after being banned", newTestCase().
		Given(
			chat.RoomWasOnBoarded{
				RoomID:   roomID,
				UserID:   userID,
				RoomName: roomName,
			},
			chat.UserWasBannedFromRoom{
				RoomID:  roomID,
				UserID:  userID,
				Reason:  "language",
				Timeout: 3600,
			},
		).
		When(chat.JoinRoom{
			RoomID: roomID,
			UserID: userID,
		}).
		Then())
}

func Test_SendMessageToRoom(t *testing.T) {
	t.Run("sends message to room", newTestCase().
		Given(
			chat.UserWasOnBoarded{
				UserID: userID,
				Name:   "John",
			},
			chat.RoomWasOnBoarded{
				RoomID:   roomID,
				UserID:   userID,
				RoomName: roomName,
			},
		).
		When(chat.SendMessageToRoom{
			RoomID:  roomID,
			UserID:  userID,
			Message: message,
		}).
		Then(chat.MessageWasSentToRoom{
			RoomID:  roomID,
			UserID:  userID,
			Message: message,
		}))

	t.Run("fails from invalid room", newTestCase().
		Given(chat.UserWasOnBoarded{
			UserID: userID,
			Name:   "John",
		}).
		When(chat.SendMessageToRoom{
			RoomID:  roomID,
			UserID:  userID,
			Message: message,
		}).
		Then())

	t.Run("bad word results in warning to user", newTestCase().
		Given(
			chat.UserWasOnBoarded{
				UserID: userID,
				Name:   "John",
			},
			chat.RoomWasOnBoarded{
				RoomID:   roomID,
				UserID:   userID,
				RoomName: roomName,
			},
		).
		When(chat.SendMessageToRoom{
			RoomID:  roomID,
			UserID:  userID,
			Message: badMessage,
		}).
		Then(
			chat.MessageWasSentToRoom{
				RoomID:  roomID,
				UserID:  userID,
				Message: badMessage,
			},
			chat.PrivateMessageWasSentToRoom{
				RoomID:       roomID,
				TargetUserID: userID,
				Message:      warnMessage,
			},
			chat.UserWasWarned{
				UserID: userID,
				Reason: "language",
			},
		))

	t.Run("2nd bad word results in 2nd warning to user", newTestCase().
		Given(
			chat.UserWasOnBoarded{
				UserID: userID,
				Name:   "John",
			},
			chat.RoomWasOnBoarded{
				RoomID:   roomID,
				UserID:   userID,
				RoomName: roomName,
			},
			chat.UserWasWarned{
				UserID: userID,
				Reason: "language",
			},
		).
		When(chat.SendMessageToRoom{
			RoomID:  roomID,
			UserID:  userID,
			Message: badMessage,
		}).
		Then(
			chat.MessageWasSentToRoom{
				RoomID:  roomID,
				UserID:  userID,
				Message: badMessage,
			},
			chat.PrivateMessageWasSentToRoom{
				RoomID:       roomID,
				TargetUserID: userID,
				Message:      warnMessage,
			},
			chat.UserWasWarned{
				UserID: userID,
				Reason: "language",
			},
		))

	t.Run("too many bad words results in user removed and banned from room", newTestCase().
		Given(
			chat.UserWasOnBoarded{
				UserID: userID,
				Name:   "John",
			},
			chat.RoomWasOnBoarded{
				RoomID:   roomID,
				UserID:   userID,
				RoomName: roomName,
			},
			chat.UserWasWarned{
				UserID: userID,
				Reason: "language",
			},
			chat.UserWasWarned{
				UserID: userID,
				Reason: "language",
			},
		).
		When(chat.SendMessageToRoom{
			RoomID:  roomID,
			UserID:  userID,
			Message: badMessage,
		}).
		Then(
			chat.MessageWasSentToRoom{
				RoomID:  roomID,
				UserID:  userID,
				Message: badMessage,
			},
			chat.UserWasRemovedFromRoom{
				RoomID: roomID,
				UserID: userID,
				Reason: "language",
			},
			chat.UserWasBannedFromRoom{
				RoomID:  roomID,
				UserID:  userID,
				Reason:  "language",
				Timeout: 3600,
			},
		))
}

func Test_SendPrivateMessageToRoom(t *testing.T) {
	t.Run("sends private message to room", newTestCase().
		Given(
			chat.UserWasOnBoarded{
				UserID: userID,
				Name:   "John",
			},
			chat.RoomWasOnBoarded{
				RoomID:   roomID,
				UserID:   userID,
				RoomName: roomName,
			},
		).
		When(chat.SendPrivateMessageToRoom{
			RoomID:       roomID,
			TargetUserID: userID,
			Message:      warnMessage,
		}).
		Then(chat.PrivateMessageWasSentToRoom{
			RoomID:       roomID,
			TargetUserID: userID,
			Message:      warnMessage,
		}))

	t.Run("fails from invalid room", newTestCase().
		Given(
			chat.UserWasOnBoarded{
				UserID: userID,
				Name:   "John",
			},
		).
		When(chat.SendPrivateMessageToRoom{
			RoomID:       roomID,
			TargetUserID: userID,
			Message:      warnMessage,
		}).
		Then())
}

func TestNew_Failures(t *testing.T) {
	t.Run("unable to subscribe", func(t *testing.T) {
		// Given
		failingStore := rangedbtest.NewFailingSubscribeEventStore()

		// When
		app, err := chat.New(failingStore)

		// Then
		require.EqualError(t, err, "failingRecordSubscription.StartFrom")
		assert.Nil(t, app)
	})
}

func newTestCase() *bdd.TestCase {
	store := inmemorystore.New()
	chat.BindEvents(store)
	app, err := chat.New(store)
	if err != nil {
		log.Fatal(err)
	}

	return bdd.New(store, func(command bdd.Command) {
		app.Dispatch(command)
		waitForProjections()
	})
}

func waitForProjections() {
	time.Sleep(10 * time.Millisecond)
}
