package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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

type fifoEmitter int

func newFifoEmitter() *fifoEmitter {
	path := filepath.Join(os.TempDir(), fifoFileName)
	syscall.Mkfifo(path, 0666)
	return new(fifoEmitter)
}

func (x *fifoEmitter) emit(evs events) error {
	path := filepath.Join(os.TempDir(), fifoFileName)
	file, err := os.OpenFile(path, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}

	var r bytes.Buffer
	for _, e := range evs {
		r.WriteString("\t- [" + e.Timestamp() + "] ")
		r.WriteString(e.Entry() + "\n")
	}

	num, err := io.Copy(file, &r)
	return err
}
