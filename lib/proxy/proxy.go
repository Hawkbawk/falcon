package proxy

import (
	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

const ProxyImageName = "Hawkbawk/falcon-proxy"
const ProxyContainerName = "falcon-proxy"

var containerConfig container.Config = container.Config{
	Image: ProxyImageName,
	ExposedPorts: nat.PortSet{
		"80":   struct{}{},
		"8080": struct{}{},
	},
}

var hostConfig container.HostConfig = container.HostConfig{
	Binds: []string{
		"/var/run/docker.sock:/var/run/docker.sock:ro",
	},
	PortBindings: nat.PortMap{
		"80": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "80",
			},
		},
		"8080": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "8080",
			},
		},
	},
}

// StartProxy starts up the falcon-proxy so that it can start forwarding requests.
func StartProxy() error {
	if err := docker.StartContainer(ProxyImageName, &hostConfig, &containerConfig, ProxyContainerName); err != nil {
		return err
	}
	return nil
}
