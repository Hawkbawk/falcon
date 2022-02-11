package docker_test

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/Hawkbawk/falcon/mocks/mock_docker"
)

func MockContainerListWithValues(containers []types.Container, err error, mockClient *mock_docker.MockDockerClient, containerName string) {
	mockClient.EXPECT().ContainerList(context.Background(),
		types.ContainerListOptions{All: true, Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName})}).Return(containers, err)
}

func MockContainerRemoveWithError(id string, err error, mockClient *mock_docker.MockDockerClient) {
	mockClient.EXPECT().ContainerRemove(context.Background(), id,
		types.ContainerRemoveOptions{Force: true}).Return(err)
}

func ExpectContainerRemoveNotBeCalled(id string, mockClient *mock_docker.MockDockerClient) {
	mockClient.EXPECT().ContainerRemove(context.Background(), id,
		types.ContainerRemoveOptions{Force: true}).Times(0)
}

var _ = Describe("Docker", func() {
	var (
		ctrl          *gomock.Controller
		mockClient    *mock_docker.MockDockerClient
		containerName = "testcontroller"
		containerId = "abcd"

	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockClient = mock_docker.NewMockDockerClient(ctrl)
	})

	Describe("GetContainerID", func() {
		Describe("the container exists", func() {
			var (
				containerList []types.Container
			)

			BeforeEach(func() {
				var c = types.Container{
					ID: containerId,
				}
				containerList = make([]types.Container, 0, 1)
				containerList = append(containerList, c)
			})

			It("returns the containers id and no errors", func() {
				MockContainerListWithValues(containerList, nil, mockClient, containerName)

				result, err := docker.GetContainer(containerName, mockClient)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.ID).To(Equal(containerId))
			})
		})

		Describe("the container doesn't exist", func() {
			It("returns no id and no errors", func() {
				MockContainerListWithValues(make([]types.Container, 0), nil, mockClient, containerName)

				result, err := docker.GetContainer(containerName, mockClient)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).Should(BeNil())
			})
		})

		Describe("the client returns an error", func() {
			It("returns no id and an error", func() {
				MockContainerListWithValues(make([]types.Container, 0), fmt.Errorf("err"), mockClient, containerName)

				result, err := docker.GetContainer(containerName, mockClient)
				Expect(err).Should(MatchError("err"))
				Expect(result).Should(BeNil())
			})
		})
	})

	Describe("RemoveContainer", func() {
		Describe("the container exists", func() {
			var (
				containerList []types.Container
			)

			BeforeEach(func() {
				var c = types.Container{
					ID: containerId,
				}
				containerList = make([]types.Container, 0, 1)
				containerList = append(containerList, c)
			})

			It("removes the container and returns no errors", func() {
				MockContainerListWithValues(containerList, nil, mockClient, containerName)
				MockContainerRemoveWithError(containerId, nil, mockClient)

				Expect(docker.RemoveContainer(containerName, mockClient)).Should(Succeed())
			})

			Describe("the container removal throws an error", func() {
				It("returns an error", func() {
					MockContainerListWithValues(containerList, nil, mockClient, containerName)
					MockContainerRemoveWithError(containerId, fmt.Errorf("err"), mockClient)

					Expect(docker.RemoveContainer(containerName, mockClient)).Should(MatchError("err"))
				})
			})
		})

		Describe("the container doesn't exist", func() {
			It("doesn't try and remove the container and doesn't error", func() {
				MockContainerListWithValues(make([]types.Container, 0), nil, mockClient, containerName)
				ExpectContainerRemoveNotBeCalled(containerId, mockClient)
				Expect(docker.RemoveContainer(containerName, mockClient)).Should(Succeed())
			})
		})

		Describe("the container find throws an error", func() {
			It("doesn't try and remove the container and returns an error", func() {
				MockContainerListWithValues(make([]types.Container, 0), fmt.Errorf("err"), mockClient, containerName)
				ExpectContainerRemoveNotBeCalled(containerId, mockClient)
				Expect(docker.RemoveContainer(containerName, mockClient)).Should(MatchError("err"))
			})
		})
	})
})
