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

func loop(ipc chan signal) {
	ticker := time.NewTicker(syncInterval)

	c, err := newCvlt(cvltConfig{
		// root:  "/data/docker/containers",
		root:    "testdata",
		match:   "*-json.log",
		depth:   1,
		logType: logTypeJSON,
	})
	if err != nil {
		ipc <- signal{
			code: sigInitError,
			err:  err}
		return
	}
	pool := []*cvlt{c}

	for {
		select {
		case <-ticker.C:
			for _, entry := range pool {
				entry.sync(ipc)
			}
		}
	}
	ipc <- signal{code: sigQuit}
}
