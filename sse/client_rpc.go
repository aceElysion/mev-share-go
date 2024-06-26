package sse

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// InternalClient is a client for the matchmaker
type InternalClient struct {
	BaseURL string // BaseURL is the base URL for the matchmaker
}

// New creates a new InternalClient for the matchmaker with the given base URL
func New(baseURL string) SSEClient {
	return &InternalClient{
		BaseURL: baseURL,
	}
}

// Subscription represents a subscription to matchmaker events
type Subscription struct {
	client    http.Client
	rspBody   io.ReadCloser
	stopper   chan struct{}
	eventChan chan<- Event
}

// Subscribe to matchmaker events and returns a type that can be used to control the subscription
func (c *InternalClient) Subscribe(eventChan chan<- Event) (SSESubscription, error) {
	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	sub := &Subscription{
		client:    client,
		rspBody:   resp.Body,
		eventChan: eventChan,
		stopper:   make(chan struct{}),
	}

	go sub.readEvents()

	return sub, nil
}

// readEvents reads the events and sends them to the event channel
func (s *Subscription) readEvents() {
	defer s.rspBody.Close()

	var err error
	var event MatchMakerEvent
	scanner := bufio.NewScanner(s.rspBody)

	for {
		// send event or stop
		select {
		case <-s.stopper:
			close(s.eventChan)
			close(s.stopper)
			return
		default:
			if err == nil {
				if scanner.Scan() {
					data := scanner.Text()
					if len(data) > 0 && data != ":ping" {
						data = strings.TrimPrefix(data, "data: ")
						_ = json.Unmarshal([]byte(data), &event)
					}
				} else {
					// 出现error
					err = scanner.Err()
					if err == nil {
						err = errors.New("EOF")
					}
				}

				s.eventChan <- Event{
					Data:  &event,
					Error: err,
				}
			}
		}
	}
}

// Stop stops the subscription to matchmaker events
func (s *Subscription) Stop() {
	s.stopper <- struct{}{}
}
