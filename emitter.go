package main

import "fmt"

type emitter interface {
	emit(events) error
}

type cliEmitter int

func (x cliEmitter) emit(evs events) error {
	for _, e := range evs {
		fmt.Printf("- %q\n", e)
	}
	return nil
}
