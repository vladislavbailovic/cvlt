package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type event fmt.Stringer
type events []event

type rawEvent string

func (x rawEvent) String() string {
	return string(x)
}

type jsonLogEvent struct {
	Time string `json:"time"`
	Log  string `json:"log"`
}

func (x jsonLogEvent) String() string {
	return fmt.Sprintf("[%s] %s", x.Time, x.Log)
}

type eventParser func([]byte) (events, error)

func parseEvents(b []byte) (events, error) {
	var result events
	if len(b) == 0 {
		return result, nil
	}

	splits := strings.Split(string(b), "\n")
	result = make(events, 0, len(splits))
	for _, s := range splits {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		result = append(result, NewEvent(s))
	}
	return result, nil
}

func NewEvent(raw string) event {
	var ev jsonLogEvent
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		return rawEvent(raw)
	}
	return ev
}

func (x events) emit(to []emitter) error {
	for _, event := range x {
		for _, e := range to {
			if err := e.update(event); err != nil {
				return err
			}
		}
	}
	return nil
}
