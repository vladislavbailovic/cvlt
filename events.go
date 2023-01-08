package main

type events []event

type event interface {
	Timestamp() string
	Entry() string
}

type rawEvent string

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

func (x jsonLogEvent) Timestamp() string {
	return x.Time
}

func (x jsonLogEvent) Entry() string {
	return x.Log
}
