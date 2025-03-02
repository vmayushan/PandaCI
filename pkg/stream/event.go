package stream

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
)

type EventStream struct {
	*Stream
}

type EventType string

const (
	EventTypeUnknown EventType = "unknown"
)

type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	Data      any       `json:"data,omitempty"`
}

func (es *EventStream) UnknownEvent() error {
	return es.writeEvent(Event{
		Timestamp: time.Now(),
		Type:      EventTypeUnknown,
	})
}

func (e *EventStream) writeEvent(event Event) error {

	data, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("marshalling event")
		return err
	}

	e.Write(string(data))
	return nil
}

func NewEventStream() *EventStream {
	return &EventStream{
		Stream: NewStream([]string{}),
	}
}
