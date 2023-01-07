package main

import (
	"os"
	"testing"
)

func Test_Update_ErrorsWithWrongFile(t *testing.T) {
	x := &cursor{path: "whatever this does not exist"}
	if err := x.update(); err == nil {
		t.Error("expected error")
	}
	if x.size != 0 || x.pos != 0 {
		t.Errorf("expected zero initialization: %#v", x)
	}
}

func Test_Update_HappyPath(t *testing.T) {
	x := &cursor{path: "testdata/test.log"}
	if err := x.update(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if x.size == 0 {
		t.Error("expected log size to not be zero")
	}
	if x.pos != 0 {
		t.Errorf("expected zero pos, got %d", x.pos)
	}
}

func Test_isChanged(t *testing.T) {
	x := &cursor{path: "testdata/test.log"}
	if x.isChanged() {
		t.Error("expected unchanged initially")
	}

	x.update()
	if !x.isChanged() {
		t.Error("expected updated after change")
	}

	x.earmark()
	if x.isChanged() {
		t.Error("expected unchanged after record")
	}
}

func Test_Latest(t *testing.T) {
	x := &cursor{path: "testdata/test.log"}
	x.update()
	got, err := x.latest()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want, err := os.ReadFile("testdata/test.log")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if string(want) != string(got) {
		t.Errorf("buffer mismatch, want %v got %v", want, got)
	}
}
