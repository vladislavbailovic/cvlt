package main

import "time"

type events []event

type event interface {
	Timestamp() time.Time
	Entry() string
}

type rawEvent string

func (x rawEvent) Timestamp() time.Time {
	return time.Now()
}

func (x rawEvent) Entry() string {
	return string(x)
}

type jsonLogEvent struct {
	Time string `json:"time"`
	Log  string `json:"log"`
}

func (x jsonLogEvent) Timestamp() time.Time {
	tm, err := time.Parse(time.RFC3339Nano, x.Time)
	if err != nil {
		return time.Now()
	}
	return tm
}

func (x jsonLogEvent) Entry() string {
	return x.Log
}
