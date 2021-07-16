package docker

import (
	"fmt"

	"github.com/hawkbawk/prox/lib/logger"
	"github.com/hawkbawk/prox/lib/shell"
)

const DefaultNetworkName = "prox_network"
const defaultGateway = "172.26.0.1"
const defaultSubnet = "172.26.0.0/16"

func CreateDefaultNetwork(networkName string) {
	if networkName == "" {
		networkName = DefaultNetworkName
	}

	logger.LogInfo("Creating Docker network needed for Traefik to work...")

	args := []string{
		"network",
		"create",
		networkName,
		fmt.Sprint("--gateway", "=", defaultGateway),
		fmt.Sprint("--subnet", "=", defaultSubnet),
	}

	err := shell.RunCommand("docker", args)
	if e, ok := err.(*shell.CommandFailed); ok {
		logger.LogError("Couldn't create the Docker network named", networkName,
			".\nPlease see the following error message: ", e.Error())
	}
}
