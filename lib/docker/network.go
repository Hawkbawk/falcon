package docker

import (
	"fmt"

	"github.com/hawkbawk/falcon/lib/logger"
	"github.com/hawkbawk/falcon/lib/shell"
)

// The gateway points to the Traefik container on the Docker network.
const DefaultGateway = "172.26.0.1"
const defaultNetworkName = "prox_network"
const defaultSubnet = "172.26.0.0/16"

var createdNetworkName string

// Creates a new Docker network that the Traefik container will listen to in order
// to act as a reverse-proxy. A custom network name can be defined. If one isn't passed,
// then a default name is used.
func CreateDockerNetwork(networkName string) {
	if networkName == "" {
		networkName = defaultNetworkName
	}

	logger.LogInfo("Creating Docker network needed for Traefik to work...")

	args := []string{
		"network",
		"create",
		networkName,
		fmt.Sprint("--gateway", "=", DefaultGateway),
		fmt.Sprint("--subnet", "=", defaultSubnet),
	}

	err := shell.RunCommand("docker", args, false)
	if err != nil {
		logger.LogError("Couldn't create the Docker network named", networkName,
			".\nPlease see the following error message for more details: ", err.Error())
	}

	createdNetworkName = networkName
}

// DestroyDockerNetwork destroys the Docker network that was previously created by
// CreateDockerNetwork command.
func DestroyDockerNetwork() {
	logger.LogInfo("Destroying the Traefik Docker network...")

	args := []string{
		"network",
		"rm",
		createdNetworkName,
	}
	err := shell.RunCommand("docker", args, false)

	if err != nil {
		logger.LogError("Unable to delete the Traefik Docker Network.\n",
			"Please see the following error for more details: ", err.Error())
	}
}
