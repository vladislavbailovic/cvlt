package main

import (
	"fmt"
	"os"
	ossig "os/signal"
	"syscall"
	"time"
)

var syncInterval time.Duration = time.Millisecond * 500
var fifoFileName string = "cvlt.fifo"

func main() {
	ipc := make(chan signal)
	go loop(ipc)
	for {
		select {
		case sig := <-ipc:
			if sig.err != nil {
				fmt.Fprintf(
					os.Stderr,
					fmt.Sprintf("[ERROR] %s: %v\n", sig.code.String(), sig.err))
			}
			if !sig.canContinue() {
				os.Exit(sig.exitCode())
			}
		}
	}
}

func loop(ipc chan signal) {
	ticker := time.NewTicker(syncInterval)

	cfgs := []cvltConfig{
		cvltConfig{
			// root:  "/data/docker/containers",
			root:    "testdata",
			match:   "*-json.log",
			depth:   1,
			logType: logTypeJSON,
		},
		cvltConfig{
			root:    "testdata",
			match:   "test*.log",
			depth:   1,
			logType: logTypePlain,
		},
	}
	pool := make([]*cvlt, 0, len(cfgs))

	for _, cfg := range cfgs {
		c, err := newCvlt(cfg)
		if err != nil {
			ipc <- signal{
				code: sigInitError,
				err:  err}
			return
		}
		pool = append(pool, c)
	}
	defer func() {
		// Flush at the end
		for _, entry := range pool {
			entry.flush()
		}
	}()

	outside := make(chan os.Signal)
	ossig.Notify(outside, os.Interrupt, syscall.SIGTERM)
outer:
	for {
		select {
		case <-ticker.C:
			for _, entry := range pool {
				entry.sync(ipc)
			}
		case <-outside:
			break outer
		}
	}
	ipc <- signal{code: sigQuit}
}
