package proxy

import (
	"fmt"
	"os"

	"github.com/Hawkbawk/falcon/mocks/mock_docker"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Proxy", func() {
	var (
		ctrl       *gomock.Controller
		mockClient *mock_docker.MockDockerClient
		hostname   = "example.com"
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockClient = mock_docker.NewMockDockerClient(ctrl)
	})

	Describe("Start", func() {
		It("tries to start the proxy container and returns no errors", func() {
			mockClient.EXPECT().StartContainer(ProxyImageName, hostConfig, containerConfig, ProxyContainerName).Return(nil)

			Expect(Start(mockClient)).To(Succeed())
		})

		It("returns an error if the container can't be started", func() {
			err := fmt.Errorf("problems!")
			mockClient.EXPECT().StartContainer(ProxyImageName, hostConfig, containerConfig, ProxyContainerName).Return(err)

			Expect(Start(mockClient)).To(Equal(err))
		})
	})

	Describe("Stop", func() {
		It("tries to stop the proxy container and returns no errors", func() {
			mockClient.EXPECT().StopAndRemoveContainer(ProxyContainerName).Return(nil)

			Expect(Stop(mockClient)).To(Succeed())
		})

		It("returns an error if the container can't be stopped", func() {
			err := fmt.Errorf("problems!")
			mockClient.EXPECT().StopAndRemoveContainer(ProxyContainerName).Return(err)

			Expect(Stop(mockClient)).To(Equal(err))
		})
	})

	Describe("createTlsFiles", func() {
		var (
			argList   []string
			err       error
			cmdRunner = func(cmd string) error {
				argList = append(argList, cmd)
				return err
			}
		)

		BeforeEach(func() {
			argList = make([]string, 0)
			err = nil
		})

		It("tries to call mkcert to make the certs in the right directory", func() {
			Expect(createTlsFiles(hostname, cmdRunner)).To(Succeed())

			currDir, err := os.Getwd()

			if err != nil {
				Fail("Unable to check current working directory to verify we're in the certificates directory.")
			}
			Expect(argList[0]).To(Equal(fmt.Sprintf("mkcert %v", hostname)))
			Expect(currDir).To(Equal(certificatesDir))
		})

		It("returns an error if the command fails", func() {
			err = fmt.Errorf("couldn't do that chief!")

			result := createTlsFiles(hostname, cmdRunner)
			Expect(result).To(Equal(err))
		})
	})

	Describe("addFilesToConfig", func() {
		var (
			testData = `
tls:
  certificates:
    - certFile: /example.pem
		  keyFile: /example-key.pem
`
			expectedResult map[interface{}]interface{}
		)

		BeforeEach(func() {

			expectedResult = make(map[interface{}]interface{})

			expectedConfig := `
tls:
  certificates:
    - certFile: /example.pem
		  keyFile: /example-key.pem
    - certFile: /usr/src/app/certs/example.com.pem
      keyFile: /usr/src/app/certs/example.com-key.pem
`
			yaml.Unmarshal([]byte(expectedConfig), &expectedResult)
		})

		It("adds the hostname files to the config", func() {
			result, _ := addFilesToConfig(hostname, []byte(testData))

			resultUnmarshalled := make(map[interface{}]interface{})

			yaml.Unmarshal(result, &resultUnmarshalled)
			Expect(resultUnmarshalled).To(Equal(expectedResult))
		})

		It("returns an error if it can't unmarshal the config", func() {
			testData = `
tls:
  certificates:
	  ---- invalid yaml -----
`
			Expect(addFilesToConfig(hostname, []byte(testData))).Error().Should(HaveOccurred())
		})
	})

	Describe("createCertFileName", func() {
		It("creates the right file name", func() {
			Expect(createCertFileName(hostname)).To(Equal("example.com.pem"))
		})
	})

	Describe("createKeyFileName", func() {
		It("creates the right file name", func() {
			Expect(createKeyFileName(hostname)).To(Equal("example.com-key.pem"))
		})
	})
})
