package dnsmasq

import (
	"fmt"

	"github.com/Hawkbawk/falcon/mocks/mock_docker"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dnsmasq", func() {
	var (
		ctrl       *gomock.Controller
		mockClient *mock_docker.MockDockerClient
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockClient = mock_docker.NewMockDockerClient(ctrl)
	})

	Describe("Start", func() {
		It("tries to start the dnsmasq container and returns no errors", func() {
			mockClient.EXPECT().StartContainer(dnsMasqImageName, hostConfig, containerConfig, dnsMasqContainerName).Return(nil)

			Expect(Start(mockClient)).Should(Succeed())
		})

		It("returns an error if the container can't be started", func() {
			err := fmt.Errorf("problems!")
			mockClient.EXPECT().StartContainer(dnsMasqImageName, hostConfig, containerConfig, dnsMasqContainerName).Return(err)

			Expect(Start(mockClient)).Should(Equal(err))
		})
	})

	Describe("Stop", func() {
		It("tries to stop the dnsmasq container and returns no errors", func() {
			mockClient.EXPECT().StopAndRemoveContainer(dnsMasqContainerName).Return(nil)

			Expect(Stop(mockClient)).Should(Succeed())
		})

		It("returns an error if the container can't be stopped", func() {
			err := fmt.Errorf("problems!")
			mockClient.EXPECT().StopAndRemoveContainer(dnsMasqContainerName).Return(err)

			Expect(Stop(mockClient)).Should(Equal(err))
		})
	})
})
