package main

type source struct {
	pos cursor
}

func newSource(path string) *source {
	return &source{
		pos: cursor{path: path},
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
