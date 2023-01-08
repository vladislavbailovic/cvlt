package main

import (
	"encoding/json"
	"strings"
)

type eventParser func([]byte) (events, error)

type logType uint

const (
	logTypePlain logType = iota
	logTypeJSON
)

func newEventParser(kind logType) eventParser {
	switch kind {
	case logTypeJSON:
		return parseJsonLogline
	}
	return parsePlaintextLogline
}

func parsePlaintextLogline(b []byte) (events, error) {
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
		result = append(result, rawEvent(s))
	}
	return result, nil
}

func parseJsonLogline(b []byte) (events, error) {
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
		var ev jsonLogEvent
		if err := json.Unmarshal([]byte(s), &ev); err != nil {
			return result, err
		}
		result = append(result, ev)
	}
	return result, nil
}
