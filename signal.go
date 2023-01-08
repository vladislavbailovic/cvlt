package main

type code int

const (
	sigQuit code = iota
	sigInitError
	sigUpdateError
	sigParseError
	sigBroadcastError
)

func (x code) String() string {
	switch x {
	case sigQuit:
		return "Quit"
	case sigInitError:
		return "Initialization Error"
	case sigUpdateError:
		return "Update Error"
	case sigParseError:
		return "Event Parsing Error"
	case sigBroadcastError:
		return "Broadcast Error"
	}
	return "Unknown Signal"
}

type signal struct {
	code code
	err  error
}

func (x signal) canContinue() bool {
	if x.code == sigBroadcastError {
		return true
	}
	return false
}

func (x signal) exitCode() int {
	return int(x.code)
}
