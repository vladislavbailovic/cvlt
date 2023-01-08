package main

import (
	"path/filepath"
	"strings"
)

type cvlt struct {
	following []*source
	parser    eventParser
	audience  []emitter
}

func newCvlt(cfg cvltConfig) (*cvlt, error) {
	var tails []*source

	for i := 0; i < cfg.depth; i++ {
		nest := strings.Repeat("*", i)
		matches, err := filepath.Glob(
			filepath.Join(cfg.root, nest, cfg.match))
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			tails = append(tails,
				newSource(match))
		}
	}

	return &cvlt{
		following: tails,
		parser:    parseEvents,
		audience:  []emitter{new(cliEmitter)},
	}, nil
}

func (x cvlt) sync(ipc chan signal) {
	for _, item := range x.following {
		go func(item *source) {
			change, err := item.sync()
			if err != nil {
				ipc <- signal{
					code: sigUpdateError,
					err:  err}
				return
			}
			err = x.broadcast(change)
			if err != nil {
				ipc <- signal{
					code: sigBroadcastError,
					err:  err}
				return
			}
		}(item)
	}
}

func (x cvlt) broadcast(change []byte) error {
	evs, err := x.parser(change)
	if err != nil {
		return err
	}
	evs.emit(x.audience)
	return nil
}

type cvltConfig struct {
	root  string
	depth int
	match string
}