package platform

import (
	"fmt"
	"strings"

	"github.com/juju/errors"

	"github.com/k8scat/articli/pkg/csdn"
	"github.com/k8scat/articli/pkg/juejin"
	"github.com/k8scat/articli/pkg/oschina"
	"github.com/k8scat/articli/pkg/segmentfault"
)

var (
	hub           = map[string]Platform{}
	platformNames []string
)

func init() {
	register(new(juejin.Client))
	register(new(csdn.Client))
	register(new(oschina.Client))
	register(new(segmentfault.Client))
}

func register(p Platform) {
	if _, reg := hub[p.Name()]; reg {
		panic(fmt.Sprintf("platform %s already registered", p.Name()))
	}
	hub[p.Name()] = p
	platformNames = append(platformNames, p.Name())
}

func GetByName(name string) (Platform, error) {
	pf, ok := hub[name]
	if ok {
		return pf, nil
	}
	return nil, errors.Errorf("platform [%s] not supported, current supported platforms: %s", name, strings.Join(platformNames, ", "))
}
