package platform

import (
	"fmt"
	"strings"

	"github.com/k8scat/articli/pkg/csdn"
	"github.com/k8scat/articli/pkg/juejin"
	"github.com/k8scat/articli/pkg/oschina"
)

var (
	hub           = map[string]Platform{}
	platformNames []string
)

func init() {
	register(new(juejin.Client))
	register(new(csdn.Client))
	register(new(oschina.Client))
}

func register(p Platform) error {
	if _, reg := hub[p.Name()]; reg {
		return fmt.Errorf("platform %s already registered", p.Name())
	}
	hub[p.Name()] = p
	platformNames = append(platformNames, p.Name())
	return nil
}

func GetByName(name string) (Platform, error) {
	pf, ok := hub[name]
	if ok {
		return pf, nil
	}
	return nil, fmt.Errorf("platform [%s] not supported, current supported platforms: %s", name, strings.Join(platformNames, ", "))
}
