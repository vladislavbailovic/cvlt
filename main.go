package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var syncInterval time.Duration = time.Millisecond * 500

func main() {
	ipc := make(chan signal)
	go loop(ipc)
	select {
	case signal := <-ipc:
		switch signal {
		case squit:
			os.Exit(0)
		}
	}
}

type signal int

const (
	squit signal = iota
	sfatal
)

type cvlt []*followed

func loop(ipc chan signal) {
	ticker := time.NewTicker(syncInterval)
	pool := cvlt{
		newFollowed("testdata/test.log"),
		newFollowed("testdata/test2.log"),
	}
	for {
		select {
		case <-ticker.C:
			for _, item := range pool {
				go func(item *followed) {
					evs, err := sync(&item.pos)
					if err != nil {
						fmt.Println("[ERROR]", err)
						ipc <- sfatal
						return
					}
					evs.emit(item.emitters)
				}(item)
			}
		}
	}
	ipc <- squit
}

func sync(item *cursor) (events, error) {
	var evs events
	if err := item.update(); err != nil {
		return evs, err
	}
	if changed, err := item.isChanged(); err != nil {
		return evs, err
	} else if !changed {
		return evs, nil
	}

	buffer, err := item.latest()
	if err != nil {
		return evs, err
	}
	item.earmark()

	return parse(buffer)
}

func parse(b []byte) (events, error) {
	splits := strings.Split(string(b), "\n")
	result := make(events, 0, len(splits))
	for _, s := range splits {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		result = append(result, s)
	}
	return result, nil
}

type events []string

func (x events) emit(to []emitter) error {
	for _, event := range x {
		for _, e := range to {
			if err := e.broadcast(event); err != nil {
				return err
			}
		}
	}
	return nil
}
