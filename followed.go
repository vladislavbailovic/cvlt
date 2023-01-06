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
