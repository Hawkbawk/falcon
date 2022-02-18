package docker

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerApi interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error
	ContainerRestart(ctx context.Context, containerID string, timeout *time.Duration) error
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
}

type DockerClient interface {
	// GetContainer finds the first container that matches the specified container name.
	// If no match is found, then a nil container and nil error is returned. Note that this function only
	// looks at containers that are in a running state.
	GetContainer(containerName string) (*types.Container, error)
	// Stops and removes the first container that matches the provided container name.
	// If no containers match, nothing happens. If any errors are encountered, they're returned.
	StopAndRemoveContainer(containerName string) error
	// Starts a container with the specified configuration. If any errors are encountered, they are returned.
	StartContainer(imageName string, hostConfig *container.HostConfig, containerConfig *container.Config, containerName string) error
}

type dockerConsumer struct {
	api DockerApi
}

func NewDockerClient() (DockerClient, error) {
	api, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)

	if err != nil {
		return nil, err
	}

	return dockerConsumer{
		api: api,
	}, nil
}

func (dc dockerConsumer) GetContainer(containerName string) (*types.Container, error) {
	ctx := context.Background()

	containers, err := dc.api.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})})

	if err != nil {
		return nil, err
	} else if len(containers) == 0 {
		return nil, nil
	} else {
		return &containers[0], nil
	}
}

func (dc dockerConsumer) StopAndRemoveContainer(containerName string) error {
	ctx := context.Background()

	container, err := dc.GetContainer(containerName)

	if err != nil {
		return err
	} else if container == nil {
		return nil
	}

	// Make sure we stop the container before removing it, cause the SDK docs lie, and it doesn't
	// actually kill the container for you, unless you specify force.
	if err := dc.api.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return err
	}

	return nil
}

func (dc dockerConsumer) StartContainer(imageName string, hostConfig *container.HostConfig, containerConfig *container.Config, containerName string) error {
	ctx := context.Background()

	container, err := dc.GetContainer(containerName)

	if err != nil {
		return err
	} else if container != (*types.Container)(nil) {
		// I've been burned in the past by not checking whether a container already exists
		// with our specified name and restarting it if it's not already running.
		if container.Status != "running" {
			if err := dc.api.ContainerRestart(ctx, container.ID, nil); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	reader, err := dc.api.ImagePull(ctx, imageName, types.ImagePullOptions{})

	if err != nil {
		return err
	}

	// We have to write the stream of data from pulling the image, otherwise we
	// won't actually pull the image.
	defer reader.Close()
	io.Copy(io.Discard, reader)

	ref, err := dc.api.ContainerCreate(ctx,
		containerConfig,
		hostConfig, &network.NetworkingConfig{},
		nil, containerName)

	if err != nil {
		return err
	}

	if err := dc.api.ContainerStart(ctx, ref.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}
