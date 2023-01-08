package main

import "testing"

func Test_CanContinue(t *testing.T) {
	suite := map[code]bool{
		sigQuit:           false,
		sigUpdateError:    false,
		sigBroadcastError: true,
	}
	for c, want := range suite {
		t.Run(c.String(), func(t *testing.T) {
			s := signal{code: c}
			got := s.canContinue()
			if want != got {
				t.Errorf("can continue: want %v, got %v",
					want, got)
			}
		})
	}
}

func Test_ExitCode(t *testing.T) {
	suite := []code{
		sigQuit,
		sigInitError,
		sigUpdateError,
		sigBroadcastError,
	}
	for want, c := range suite {
		t.Run(c.String(), func(t *testing.T) {
			s := signal{code: c}
			got := s.exitCode()
			if want != got {
				t.Errorf("exit code: want %v, got %v",
					want, got)
			}
		})
	}
}
