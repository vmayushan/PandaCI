package stream

import (
	"sync"
)

type Stream struct {
	mu            sync.Mutex
	entries       []string
	subscriptions []chan []string
}

func (s *Stream) Write(data string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = append(s.entries, data)

	for _, sub := range s.subscriptions {
		sub <- []string{data}
	}
}

func (s *Stream) Entries() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	safeArray := make([]string, len(s.entries))

	copy(safeArray, s.entries)

	return safeArray
}

func (s *Stream) Subscribe() chan []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub := make(chan []string, 100)
	s.subscriptions = append(s.subscriptions, sub)

	sub <- s.entries

	return sub
}

func (s *Stream) Unsubscribe(sub chan []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, subscription := range s.subscriptions {
		if subscription == sub {
			s.subscriptions = append(s.subscriptions[:i], s.subscriptions[i+1:]...)
			close(sub)
			return
		}
	}
}

func NewStream(inital []string) *Stream {
	return &Stream{
		entries: inital,
	}
}
