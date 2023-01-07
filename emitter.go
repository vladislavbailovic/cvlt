package main

import "fmt"

type emitter interface {
	update(event) error
}

type cliEmitter int

func (x cliEmitter) update(e event) error {
	fmt.Printf("- %q\n", e)
	return nil
}
