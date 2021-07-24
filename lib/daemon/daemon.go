package daemon

import (
	"github.com/kardianos/service"
)

const Name = "falcon-daemon"
const DisplayName = "falcon-daemon"
const Description = `Keeps the falcon-proxy container in sync with Docker networks on this machine.`

var svcConfig *service.Config = &service.Config{
	Name:        Name,
	DisplayName: DisplayName,
	Description: Description,
}

type Daemon struct {
	Syncer  ProxySyncer
	Service service.Service
}

func NewDaemon() (*Daemon, error) {
	syncer, err := NewProxySyncer()

	if err != nil {
		return nil, err
	}

	service, err := service.New(syncer, svcConfig)
	if err != nil {
		return nil, err
	}

	return &Daemon{
		Syncer:  syncer,
		Service: service,
	}, nil
}
