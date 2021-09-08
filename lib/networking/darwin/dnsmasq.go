package darwin

import (
	"fmt"

	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// The name of the dnsmasq image we use.
const DnsMasqImageName = "4km3/dnsmasq:2.85-r2"
// The name of the dnsmasq container when it's running.
const DnsMasqContainerName = "falcon-dnsmasq"

var containerConfig container.Config = container.Config{
	Image: DnsMasqImageName,
	ExposedPorts: nat.PortSet{
		"53/tcp": struct{}{},
		"53/udp": struct{}{},
	},
	Cmd: []string{"-S", fmt.Sprintf("/docker/%v", loopbackAddress)}, // Tells dnsmasq to forward all requests for *.docker domains to our special loopback address.
}

var hostConfig container.HostConfig = container.HostConfig{
	PortBindings: nat.PortMap{
		"53/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "53",
			},
		},
		"53/udp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "53",
			},
		},
	},
	CapAdd: []string{"NET_ADMIN"},
}

// Starts our dnsmasq container.
func startDnsmasq() error {
	logger.LogInfo("Starting the dnsmasq container...")
	return docker.StartContainer(DnsMasqImageName, &hostConfig, &containerConfig, DnsMasqContainerName)
}

// Stops our dnsmasq container.
func stopDnsmasq() error {
	logger.LogInfo("Stopping the dnsmasq container...")
	return docker.StopContainer(DnsMasqContainerName)
}
