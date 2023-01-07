package main

type followed struct {
	pos      cursor
	emitters []emitter
}

func newFollowed(path string) *followed {
	return &followed{
		pos: cursor{path: path},
		emitters: []emitter{
			new(cliEmitter),
		},
	}
}

func (x *followed) sync() ([]byte, error) {
	var buffer []byte
	if err := x.pos.update(); err != nil {
		return buffer, err
	}
	if !x.pos.isChanged() {
		return buffer, nil
	}

	buffer, err := x.pos.latest()
	if err != nil {
		return buffer, err
	}
	x.pos.earmark()
	return buffer, nil
}

func (x *followed) broadcast(change []byte) error {
	evs, err := parseEvents(change)
	if err != nil {
		return err
	}
	evs.emit(x.emitters)
	return nil
}
