package main

import (
	"fmt"
)

type events []event

type event interface {
	fmt.Stringer
	Timestamp() string
	Entry() string
}

type rawEvent string

func (x rawEvent) String() string {
	return string(x)
}

func (x rawEvent) Timestamp() string {
	return ""
}

func (x rawEvent) Entry() string {
	return string(x)
}

type jsonLogEvent struct {
	Time string `json:"time"`
	Log  string `json:"log"`
}

func (x jsonLogEvent) String() string {
	return fmt.Sprintf("[%s] %s", x.Time, x.Log)
}

func (x jsonLogEvent) Timestamp() string {
	return x.Time
}

func (x jsonLogEvent) Entry() string {
	return x.Log
}
