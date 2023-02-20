package main

type source struct {
	pos cursor
}

func newSource(path string) *source {
	pos := cursor{path: path}
	// Update source position on spawn
	// Effectively, only sync changes *from this point on*
	// As opposed to syncing all the changes from source top
	// TODO: perhaps expose this as an option
	pos.update()
	pos.earmark()
	return &source{
		pos: pos,
	}
}

func (x *source) sync() ([]byte, error) {
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
