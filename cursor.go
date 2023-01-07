package main

import (
	"fmt"
	"os"
)

type cursor struct {
	path string
	pos  int64
	size int64
}

func (x *cursor) update() error {
	stat, err := os.Stat(x.path)
	if err != nil {
		return err
	}
	size := stat.Size()
	x.size = size
	return nil
}

func (x *cursor) isChanged() bool {
	if x.size > x.pos {
		return true
	}

	// Rotated
	if x.size < x.pos {
		x.pos = 0
		if x.size > 0 {
			// Rotated, and some stuff added
			return true
		}
	}

	// Not changed
	return false
}

func (x *cursor) latest() ([]byte, error) {
	var buffer []byte

	fp, err := os.Open(x.path)
	if err != nil {
		return buffer, nil
	}
	defer fp.Close()

	pos, err := fp.Seek(x.pos, 0)
	if err != nil {
		return buffer, nil
	}

	if pos > x.size {
		return buffer, fmt.Errorf("unreachable: current position greater than size")
	}
	buffer = make([]byte, x.size-pos)

	s, err := fp.Read(buffer)
	if s != len(buffer) {
		return buffer, fmt.Errorf("unreachanle: read length and buffer size mismatch")
	}

	return buffer, nil
}

func (x *cursor) earmark() {
	x.pos = x.size
	fmt.Println("at position", x.pos)
}
