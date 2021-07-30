package proxy

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const ProxyImageName = "hawkbawk/falcon-proxy"

var containerConfig container.Config = container.Config{
	Image: ProxyImageName,
	ExposedPorts: nat.PortSet{
		"80": struct{}{},
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
				HostIP: "0.0.0.0",
				HostPort: "80",
			},
		},
		"8080": []nat.PortBinding{
			{
				HostIP: "0.0.0.0",
				HostPort: "8080",
			},
		},
	},
}

func StartProxy() error {
	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		return err
	}
	context := context.Background()

	reader, err := client.ImagePull(context, ProxyImageName, types.ImagePullOptions{})

	if err != nil {
		return err
	}
	defer reader.Close()
	ref, err := client.ContainerCreate(context,
		&containerConfig,
		&hostConfig, &network.NetworkingConfig{},
		nil, "falcon-proxy")

	if err != nil {
		return err
	}

	if err := client.ContainerStart(context, ref.ID, types.ContainerStartOptions{}); err != nil {
		return nil
	}
	return nil
}
