package platform

import (
	"github.com/juju/errors"
)

var (
	hub = map[string]Platform{}
)

func init() {
	register(new(CSDN))
	register(new(Juejin))
}

func register(p Platform) error {
	if _, reg := hub[p.Name()]; reg {
		return errors.Errorf("Platform %s already registered", p.Name())
	}
	hub[p.Name()] = p
	return nil
}

func GetByName(name string) (Platform, bool) {
	pf, ok := hub[name]
	return pf, ok
}
