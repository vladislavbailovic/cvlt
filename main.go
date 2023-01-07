package main

import (
	"fmt"
	"os"
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
					change, err := item.sync()
					if err != nil {
						fmt.Println("[ERROR]", err)
						ipc <- sfatal
						return
					}
					err = item.broadcast(change)
					if err != nil {
						fmt.Println("[ERROR]", err)
						ipc <- sfatal
						return
					}
				}(item)
			}
		}
	}
	ipc <- squit
}
