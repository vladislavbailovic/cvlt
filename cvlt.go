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

	infix := "*"
	if cfg.depth > 1 {
		infix = "**"
	}
	name := filepath.Join(cfg.root, infix, cfg.match)

	return &cvlt{
		following: tails,
		parser:    newEventParser(cfg.logType),
		audience:  []emitter{&cliEmitter{name: name}},
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
	if len(change) == 0 {
		return nil
	}
	evs, err := x.parser(change)
	if err != nil {
		return err
	}
	for _, rcv := range x.audience {
		rcv.emit(evs)
	}
	return nil
}

type cvltConfig struct {
	root    string
	depth   int
	match   string
	logType logType
}
