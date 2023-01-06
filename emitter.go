package main

import "fmt"

type emitter interface {
	broadcast(event) error
}

type cliEmitter int

func (x cliEmitter) broadcast(e event) error {
	fmt.Printf("- %q\n", e)
	return nil
}
