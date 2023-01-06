package main

import "fmt"

type emitter interface {
	broadcast(string) error
}

type cliEmitter int

func (x cliEmitter) broadcast(e string) error {
	fmt.Printf("- %q\n", e)
	return nil
}
