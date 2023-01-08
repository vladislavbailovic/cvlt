package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func discover(dir string) cvlt {
	matches, err := filepath.Glob(filepath.Join(dir, "*", "*-json.log"))
	tails := make(cvlt, 0, len(matches))
	if err != nil {
		return tails
	}
	for _, match := range matches {
		tails = append(tails,
			newFollowed(match))
	}
	return tails
}

func loop(ipc chan signal) {
	ticker := time.NewTicker(syncInterval)
	/*
		pool := cvlt{
			// https://stackoverflow.com/a/72307042
			newFollowed("/mnt/docker/data/docker/containers/6c72bfb496a316eedaada521a17c99425cfe6f53558a8c9265200c08cf8bb6bb/6c72bfb496a316eedaada521a17c99425cfe6f53558a8c9265200c08cf8bb6bb-json.log"),
			newFollowed("testdata/test.log"),
			newFollowed("testdata/test2.log"),
		}
	*/
	pool := discover("/mnt/docker/data/docker/containers")
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
