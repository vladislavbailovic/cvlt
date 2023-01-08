package main

import "fmt"

type emitter interface {
	emit(events) error
}

type cliEmitter struct {
	name string
}

func (x *cliEmitter) emit(evs events) error {
	fmt.Println("==", x.name, "==")
	for _, e := range evs {
		fmt.Printf("\t- [%s] %q\n", e.Timestamp(), e.Entry())
	}
	return nil
}
