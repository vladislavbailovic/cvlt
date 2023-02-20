package main

import (
	"testing"
	"time"
)

func Test_jsonLogEvent_Timestamp(t *testing.T) {
	c := cursor{path: "testdata/docker-sample-1-json.log"}
	c.update()
	latest, err := c.latest()
	if err != nil {
		t.Error(err)
	}

	parser := newEventParser(logTypeJSON)
	evs, err := parser(latest)
	if err != nil {
		t.Error(err)
	}

	pivot, err := time.Parse(time.RFC3339Nano, "2023-01-08T02:00:00Z")
	if err != nil {
		t.Error(err)
	}

	for _, ev := range evs {
		ts := ev.Timestamp()
		t.Run(ts.Format(time.RFC3339Nano), func(t *testing.T) {
			if !pivot.Before(ts) {
				t.Errorf("%v should be before %v",
					pivot, ts)
			}
			if !pivot.Add(time.Hour).After(ts) {
				t.Errorf("%v should be after %v",
					pivot.Add(time.Hour), ts)
			}
		})
	}
}
