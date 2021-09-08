package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// ContainerIsRunning determines whether or not the specified container is running.
// If a non-nil error is returned, then it was unable to be determined if the specified container
// exists. Ensure you handle the error appropriately,
// otherwise you could run into undefined and unexpected behaviour.
func ContainerIsRunning(containerName string) (bool, error) {
	id, err := getContainerID(containerName)

	if err != nil {
		return false, err
	} else if id == "" {
		return false, nil
	} else {
		return true, nil
	}
}

// getContainerID determines the id of the first container that matches the specified container name.
// If no match is found, then an empty id and nil error is returned. Note that this function only
// looks at containers that are in a running state.
func getContainerID(containerName string) (string, error) {
	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName}, filters.KeyValuePair{Key: "status", Value: "running"})})

	if err != nil {
		return "", err
	} else if len(containers) == 0 {
		return "", nil
	} else {
		return containers[0].ID, nil
	}
}

func RemoveContainer(containerName string) error {
	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		return err
	}
	ctx := context.Background()

	id, err := getContainerID(containerName)

	if err != nil {
		return err
	} else if id == "" {
		return nil
	}

	if err := client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}

func StartContainer(imageName string, hostConfig *container.HostConfig, containerConfig *container.Config, containerName string) error {
	running, err := ContainerIsRunning(containerName)

	if err != nil {
		return err
	} else if running {
		return nil
	}

	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)

	if err != nil {
		return err
	}

	ctx := context.Background()

	reader, err := client.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		return err
	}
	defer reader.Close()
	ref, err := client.ContainerCreate(ctx,
		containerConfig,
		hostConfig, &network.NetworkingConfig{},
		nil, containerName)

	if err != nil {
		return err
	}

	if err := client.ContainerStart(ctx, ref.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}
