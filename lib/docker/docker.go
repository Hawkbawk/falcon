package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// getContainerID determines the id of the first container that matches the specified container name.
// If no match is found, then an empty id and nil error is returned. Note that this function only
// looks at containers that are in a running state.
func getContainer(containerName string) (*types.Container, error) {
	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	containers, err := client.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})})

	if err != nil {
		return nil, err
	} else if len(containers) == 0 {
		return nil, nil
	} else {
		return &containers[0], nil
	}
}

func RemoveContainer(containerName string) error {
	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		return err
	}
	ctx := context.Background()

	container, err := getContainer(containerName)

	if err != nil {
		return err
	} else if container == nil {
		return nil
	}

	// Make sure we stop the container before removing it, cause the SDK docs lie, and it doesn't
	// actually kill the container for you, unless you specify force.
	if err := client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return err
	}

	return nil
}

func StartContainer(imageName string, hostConfig *container.HostConfig, containerConfig *container.Config, containerName string) error {
	client, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)

	if err != nil {
		return err
	}

	ctx := context.Background()

	container, err := getContainer(containerName)

	if err != nil {
		return err
	} else if container != (*types.Container)(nil) {
		// I've been burned in the past by not checking whether a container already exists
		// with our specified name and restarting it if it's not already running.
		if container.Status != "running" {
			if err := client.ContainerRestart(ctx, container.ID, nil); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	reader, err := client.ImagePull(ctx, imageName, types.ImagePullOptions{Platform: "linux/arm64"})

	if err != nil {
		return err
	}

	// We have to write the stream of data from pulling the image, otherwise we
	// won't actually pull the image.
	defer reader.Close()
	io.Copy(io.Discard, reader)

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
