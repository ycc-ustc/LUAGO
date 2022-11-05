package state

import "luago/binchunk"
import . "luago/api"

type closure struct {
	proto *binchunk.Prototype // lua closure
	goFunc GoFunction         // go closure
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f GoFunction) *closure {
	return &closure{goFunc: f}
}
