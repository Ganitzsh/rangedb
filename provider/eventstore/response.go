package eventstore

import (
	"fmt"
	"time"
)

type NewEvent struct {
	EventID   string `json:"eventId"`
	EventType string `json:"eventType"`
	Data      string `json:"data"`
	Metadata  string `json:"metadata"`
}

type StreamResponse struct {
	Title    string    `json:"title"`
	ID       string    `json:"id"`
	Updated  time.Time `json:"updated"`
	StreamID string    `json:"streamId"`
	Author   struct {
		Name string `json:"name"`
	} `json:"author"`
	HeadOfStream bool `json:"headOfStream"`
	Links        []struct {
		URI      string `json:"uri"`
		Relation string `json:"relation"`
	} `json:"links"`
	Entries []struct {
		EventID             string    `json:"eventId"`
		EventType           string    `json:"eventType"`
		EventNumber         int       `json:"eventNumber"`
		Data                string    `json:"data"`
		StreamID            string    `json:"streamId"`
		IsJSON              bool      `json:"isJson"`
		IsMetaData          bool      `json:"isMetaData"`
		IsLinkMetaData      bool      `json:"isLinkMetaData"`
		PositionEventNumber int       `json:"positionEventNumber"`
		PositionStreamID    string    `json:"positionStreamId"`
		Title               string    `json:"title"`
		ID                  string    `json:"id"`
		Updated             time.Time `json:"updated"`
		Author              struct {
			Name string `json:"name"`
		} `json:"author"`
		Summary string `json:"summary"`
		Links   []struct {
			URI      string `json:"uri"`
			Relation string `json:"relation"`
		} `json:"links"`
	} `json:"entries"`
}

func (r *StreamResponse) NextLink() (string, error) {
	for _, link := range r.Links {
		if link.Relation == "previous" {
			return link.URI, nil
		}
	}

	return "", fmt.Errorf("next link not found")
}
