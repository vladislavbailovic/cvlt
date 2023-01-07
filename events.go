package main

import "strings"

type event string
type events []event

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
		result = append(result, event(s))
	}
	return result, nil
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
