package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Hawkbawk/falcon/mocks/mock_docker"
)

func MockContainerListWithValues(containers []types.Container, err error, mockClient *mock_docker.MockDockerApi, containerName string) {
	mockClient.EXPECT().ContainerList(context.Background(),
		types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containers, err)
}

func MockContainerRemoveWithError(id string, err error, mockClient *mock_docker.MockDockerApi) {
	mockClient.EXPECT().ContainerRemove(context.Background(), id,
		types.ContainerRemoveOptions{Force: true}).Return(err)
}

func ExpectContainerRemoveNotBeCalled(id string, mockClient *mock_docker.MockDockerApi) {
	mockClient.EXPECT().ContainerRemove(context.Background(), id,
		types.ContainerRemoveOptions{Force: true}).Times(0)
}

var _ = Describe("Docker", func() {
	var (
		ctrl          *gomock.Controller
		mockApi       *mock_docker.MockDockerApi
		containerName = "testcontroller"
		containerId   = "abcd"
		client        DockerClient
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockApi = mock_docker.NewMockDockerApi(ctrl)
		client = dockerConsumer{
			api: mockApi,
		}
	})

	Describe("GetContainer", func() {
		Describe("the container exists", func() {
			var containerList []types.Container

			BeforeEach(func() {
				containerList = []types.Container{{ID: containerId}}
			})

			It("returns the containers id and no errors", func() {
				MockContainerListWithValues(containerList, nil, mockApi, containerName)

				result, err := client.GetContainer(containerName)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.ID).To(Equal(containerId))
			})
		})

		Describe("the container doesn't exist", func() {
			It("returns no id and no errors", func() {
				MockContainerListWithValues(make([]types.Container, 0), nil, mockApi, containerName)
				Expect(client.GetContainer(containerName)).Should(BeNil())
			})
		})

		Describe("the client returns an error", func() {
			It("returns no id and an error", func() {
				MockContainerListWithValues(make([]types.Container, 0), fmt.Errorf("err"), mockApi, containerName)

				result, err := client.GetContainer(containerName)
				Expect(err).Should(MatchError("err"))
				Expect(result).Should(BeNil())
			})
		})
	})

	Describe("StopAndRemoveContainer", func() {
		Describe("the container exists", func() {
			var containerList []types.Container

			BeforeEach(func() {
				containerList = []types.Container{{ID: containerId}}
			})

			It("removes the container and returns no errors", func() {
				MockContainerListWithValues(containerList, nil, mockApi, containerName)
				MockContainerRemoveWithError(containerId, nil, mockApi)

				Expect(client.StopAndRemoveContainer(containerName)).Should(Succeed())
			})

			Describe("the container removal throws an error", func() {
				It("returns an error", func() {
					MockContainerListWithValues(containerList, nil, mockApi, containerName)
					MockContainerRemoveWithError(containerId, fmt.Errorf("err"), mockApi)

					Expect(client.StopAndRemoveContainer(containerName)).Should(MatchError("err"))
				})
			})
		})

		Describe("the container doesn't exist", func() {
			It("doesn't try and remove the container and doesn't error", func() {
				MockContainerListWithValues(make([]types.Container, 0), nil, mockApi, containerName)
				ExpectContainerRemoveNotBeCalled(containerId, mockApi)
				Expect(client.StopAndRemoveContainer(containerName)).Should(Succeed())
			})
		})

		Describe("the container find throws an error", func() {
			It("doesn't try and remove the container and returns an error", func() {
				MockContainerListWithValues(make([]types.Container, 0), fmt.Errorf("err"), mockApi, containerName)
				ExpectContainerRemoveNotBeCalled(containerId, mockApi)
				Expect(client.StopAndRemoveContainer(containerName)).Should(MatchError("err"))
			})
		})
	})
	Describe("StartContainer", func() {
		var (
			containerList []types.Container
			imageName     = "imageName"
		)

		Describe("the specified container isn't running", func() {
			BeforeEach(func() {
				containerList = []types.Container{{ID: containerId, Status: "stopped"}}
			})

			It("tries to restart the container", func() {
				mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containerList, nil)
				mockApi.EXPECT().ContainerRestart(context.Background(), containerId, nil).Return(nil)

				Expect(client.StartContainer(imageName, &container.HostConfig{}, &container.Config{}, containerName)).Should(Succeed())
			})

			It("returns the error any errors it encounters when restarting the container", func() {
				err := fmt.Errorf("problems!")
				mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containerList, nil)
				mockApi.EXPECT().ContainerRestart(context.Background(), containerId, nil).Return(err)

				Expect(client.StartContainer(imageName, &container.HostConfig{}, &container.Config{}, containerName)).Should(Equal(err))
			})
		})

		Describe("the specified container is running", func() {
			BeforeEach(func() {
				containerList = []types.Container{{ID: containerId, Status: "running"}}
			})

			It("doesn't try and do anything", func() {
				mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containerList, nil)
				mockApi.EXPECT().ImagePull(context.Background(), imageName, types.ImagePullOptions{}).Times(0)

				Expect(client.StartContainer(imageName, &container.HostConfig{}, &container.Config{}, containerName)).Should(Succeed())
			})
		})

		Describe("the specified container doesn't exist", func() {
			var (
				readCloser       io.ReadCloser
				containerConfig  = &container.Config{}
				hostConfig       = &container.HostConfig{}
				networkingConfig = &network.NetworkingConfig{}
			)

			BeforeEach(func() {
				containerList = []types.Container{}
				readCloser = io.NopCloser(strings.NewReader("testing"))
			})

			It("tries to pull the image, create the container, and then start it", func() {
				mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containerList, nil)
				mockApi.EXPECT().ImagePull(context.Background(), imageName, types.ImagePullOptions{}).Return(readCloser, nil)
				mockApi.EXPECT().ContainerCreate(context.Background(), containerConfig, hostConfig, networkingConfig, nil, containerName).Return(container.ContainerCreateCreatedBody{ID: containerId}, nil)
				mockApi.EXPECT().ContainerStart(context.Background(), containerId, types.ContainerStartOptions{}).Return(nil)

				Expect(client.StartContainer(imageName, hostConfig, containerConfig, containerName)).Should(Succeed())
			})

			Describe("error conditions", func () {
				err := fmt.Errorf("problems!")

				It("returns an error if it can't list containers", func() {
					mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(nil, err)

					Expect(client.StartContainer(imageName, hostConfig, containerConfig, containerName)).Should(Equal(err))
				})

				It("returns an error if it can't pull the image", func() {
					mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containerList, nil)
					mockApi.EXPECT().ImagePull(context.Background(), imageName, types.ImagePullOptions{}).Return(readCloser, err)

					Expect(client.StartContainer(imageName, hostConfig, containerConfig, containerName)).Should(Equal(err))
				})

				It("returns an error if it can't create the container", func() {
					mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containerList, nil)
					mockApi.EXPECT().ImagePull(context.Background(), imageName, types.ImagePullOptions{}).Return(readCloser, nil)
					mockApi.EXPECT().ContainerCreate(context.Background(), containerConfig, hostConfig, networkingConfig, nil, containerName).Return(container.ContainerCreateCreatedBody{}, err)

					Expect(client.StartContainer(imageName, hostConfig, containerConfig, containerName)).Should(Equal(err))
				})

				It("returns an error if it can't start the container", func() {
					mockApi.EXPECT().ContainerList(context.Background(), types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containerList, nil)
					mockApi.EXPECT().ImagePull(context.Background(), imageName, types.ImagePullOptions{}).Return(readCloser, nil)
					mockApi.EXPECT().ContainerCreate(context.Background(), containerConfig, hostConfig, networkingConfig, nil, containerName).Return(container.ContainerCreateCreatedBody{ID: containerId}, nil)
					mockApi.EXPECT().ContainerStart(context.Background(), containerId, types.ContainerStartOptions{}).Return(err)

					Expect(client.StartContainer(imageName, hostConfig, containerConfig, containerName)).Should(Equal(err))
				})
			})
		})
	})
})
