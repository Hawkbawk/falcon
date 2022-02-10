package dnsmasq

import (
	"fmt"

	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// The loopback address we use that makes us better than dory. It allows inter-container
// communication to work by causing all request to *.docker to resolve to 192.168.40.1, which
// containers send back out through the host networking to the falcon-proxy container, rather than
// just sending all traffic to *.docker domains back to themselves.
const LoopbackAddress = "192.168.40.1"
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
	Cmd: []string{
		"--log-facility=-", "--listen-address=0.0.0.0",
		"--interface=eth0", "--interface=docker0",
		"-A", fmt.Sprintf("/docker/%v", LoopbackAddress)}, // Tells dnsmasq to forward all requests for *.docker domains to our special loopback address.
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
func Start(client docker.DockerClient) error {
	return docker.StartContainer(DnsMasqImageName, &hostConfig, &containerConfig, DnsMasqContainerName, client)
}

// Stops our dnsmasq container.
func Stop(client docker.DockerClient) error {
	return docker.RemoveContainer(DnsMasqContainerName, client)
}
