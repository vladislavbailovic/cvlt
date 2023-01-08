package main

import (
	"fmt"
	"strings"
)

type emitter interface {
	emit(events) error
}

type cliEmitter struct {
	name string
}

func (x *cliEmitter) emit(evs events) error {
	var buf strings.Builder
	buf.WriteString("==" + x.name + "==\n")
	for _, e := range evs {
		buf.WriteString("\t- [" + e.Timestamp() + "] ")
		buf.WriteString(e.Entry() + "\n")
	}
	fmt.Println(buf.String())
	return nil
}
