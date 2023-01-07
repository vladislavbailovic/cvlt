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
	for {
		select {
		case sig := <-ipc:
			if sig.err != nil {
				fmt.Fprintf(
					os.Stderr,
					fmt.Sprintf("[ERROR] %v\n", sig.err))
			}
			if !sig.canContinue() {
				os.Exit(sig.exitCode())
			}
		}
	}
}

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
						ipc <- signal{
							code: sigUpdateError,
							err:  err}
						return
					}
					err = item.broadcast(change)
					if err != nil {
						ipc <- signal{
							code: sigBroadcastError,
							err:  err}
						return
					}
				}(item)
			}
		}
	}
	ipc <- signal{code: sigQuit}
}
