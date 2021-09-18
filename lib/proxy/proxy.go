package proxy

import (
	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

const ProxyImageName = "hawkbawk/falcon-proxy"
const ProxyContainerName = "falcon-proxy"

var containerConfig container.Config = container.Config{
	Image: ProxyImageName,
	ExposedPorts: nat.PortSet{
		"80":   struct{}{},
	},
	Labels: map[string]string{
		"traefik.enable": "true",
		"traefik.http.routers.traefik.rule": "Host(`traefik.docker`)",
		"traefik.http.services.traefik.loadbalancer.server.port": "8080",
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
	},
}

// Start starts up the falcon-proxy so that it can start forwarding requests.
func Start() {
	logger.LogInfo("Starting the falcon proxy container...")
	if err := docker.StartContainer(ProxyImageName, &hostConfig, &containerConfig, ProxyContainerName); err != nil {
		logger.LogError("Unable to start proxy container due to the following error: \n%v", err)
	}
}

func Stop() {
	logger.LogInfo("Stopping the falcon proxy container...")
	if err := docker.RemoveContainer(ProxyContainerName); err != nil {
		logger.LogError("Unable to remove the proxy container due to the following error: \n%v", err)
	}
}
