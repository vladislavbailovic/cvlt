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
	flush() error
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

func (x *cliEmitter) flush() error { return nil }

type fifoEmitter struct {
	w io.Writer
}

func newFifoEmitter() *fifoEmitter {
	path := filepath.Join(os.TempDir(), fifoFileName)
	syscall.Mkfifo(path, 0666)
	w, err := os.OpenFile(path, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	return &fifoEmitter{w: w}
}

func (x *fifoEmitter) emit(evs events) error {
	var r bytes.Buffer
	for _, e := range evs {
		r.WriteString("\t- [" + e.Timestamp() + "] ")
		r.WriteString(e.Entry() + "\n")
	}

	_, err := io.Copy(x.w, &r)
	return err
}

func (x *fifoEmitter) flush() error { return nil }
