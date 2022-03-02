package proxy

import (
	"fmt"
	"os"

	"github.com/Hawkbawk/falcon/lib/docker"
	"github.com/Hawkbawk/falcon/lib/shell"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v2"
)

const proxyImageName = "hawkbawk/falcon-proxy"
const proxyContainerName = "falcon-proxy"
const proxyConfigDir = "/usr/src/app/config"
const defaultConfig = `
# This is where falcon will add any info about any certificates that it creates for you.
# Alternatively, you can put any info about your own certificates here.
# See https://doc.traefik.io/traefik/https/tls/#user-defined for more info.
tls:
  certificates:
`

var configDir = fmt.Sprintf("%v/.falcon", os.Getenv("HOME"))
var certificatesDir = fmt.Sprintf("%v/certs", configDir)
var dynamicConfigPath = fmt.Sprintf("%v/dynamic.yml", configDir)

var containerConfig *container.Config = &container.Config{
	Image: proxyImageName,
	ExposedPorts: nat.PortSet{
		"80":  struct{}{},
		"443": struct{}{},
	},
	Labels: map[string]string{
		"traefik.enable":                                         "true",
		"traefik.http.routers.traefik.rule":                      "Host(`traefik.docker`)",
		"traefik.http.services.traefik.loadbalancer.server.port": "8080",
	},
}

var hostConfig *container.HostConfig = &container.HostConfig{
	Binds: []string{
		"/var/run/docker.sock:/var/run/docker.sock:ro",
		// We have to bind-mount the entire config directory instead of just the
		// certificates and dynamic config due to some fsnotify issues that happen
		// when you mount specific files. This ensures Traefik picks up on our
		// changes to the dynamic config.
		fmt.Sprintf("%v:%v", configDir, proxyConfigDir),
	},
	PortBindings: nat.PortMap{
		"80": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "80",
			},
		},
		"443": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "443",
			},
		},
	},
}

type TlsFilesConfig struct {
	CertFile string `yaml:"certFile,omitempty"`
	KeyFile  string `yaml:"keyFile,omitempty"`
}

type DynamicConfig struct {
	Tls struct {
		Certificates []TlsFilesConfig `yaml:"certificates,omitempty"`
	} `yaml:"tls,omitempty"`
}

// Start starts up the falcon-proxy so that it can start forwarding requests.
func Start(client docker.DockerClient) error {
	return client.StartContainer(proxyImageName, hostConfig, containerConfig, proxyContainerName)
}

// Stop stops the falcon-proxy container.
func Stop(client docker.DockerClient) error {
	return client.StopAndRemoveContainer(proxyContainerName)
}

// EnableTlsForHost creates the certificate files necessary for the specified
// hostname in the falcon certs directory and adds them to the Traefik dynamic
// config that gets mounted inside the falcon-proxy container.
func EnableTlsForHost(hostname string) error {
	if err := ensureTlsConfig(); err != nil {
		return err
	}

	// In an ideal world, we'd just be able to tell mkcert where
	// we want the certificates, but due to a bug in mkcerts arg parsing code,
	// we can't do that. That would be much easier to test, but this works for
	// now.
	if err := os.Chdir(certificatesDir); err != nil {
		return err
	}

	if err := createTlsFiles(hostname, shell.RunCommand); err != nil {
		return err
	}

	config, err := os.ReadFile(dynamicConfigPath)

	if err != nil {
		return err
	}

	newConfig, err := addFilesToConfig(hostname, config)

	if err != nil {
		return err
	}

	if err := os.WriteFile(dynamicConfigPath, []byte(newConfig), 0755); err != nil {
		return err
	}

	return nil
}

// createTlsFiles runs mkcert through the given command runner in order
// to create certificate files for the given hostname. If the given
// command runner actually runs points to a shell, then the files will
// be created in the current working directory.
func createTlsFiles(hostname string, cmdRunner func(string) error) error {
	if err := cmdRunner(fmt.Sprintf("mkcert %v", hostname)); err != nil {
		return err
	}

	return nil
}

// addFilesToConfig adds the key and cert files for the specified hostname
// to the Traefik dynamic config.
func addFilesToConfig(hostname string, data []byte) ([]byte, error) {
	config := DynamicConfig{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	certFilePath := fmt.Sprintf("%v/certs/%v", proxyConfigDir, createCertFileName(hostname))
	keyFilePath := fmt.Sprintf("%v/certs/%v", proxyConfigDir, createKeyFileName(hostname))

	config.Tls.Certificates = append(config.Tls.Certificates, TlsFilesConfig{CertFile: certFilePath, KeyFile: keyFilePath})

	return yaml.Marshal(&config)
}

// createCertFileName returns the filename for a certificate for the specified
// hostname.
func createCertFileName(hostname string) string {
	return fmt.Sprintf("%v.pem", hostname)
}

// createKeyFileName returns the filename for a key file for the specified
// hostname.
func createKeyFileName(hostname string) string {
	return fmt.Sprintf("%v-key.pem", hostname)
}

// ensureTlsConfig ensures both that the certificates directory exists,
// and that the Traefik dynamic config exists as well.
func ensureTlsConfig() error {
	if err := os.MkdirAll(certificatesDir, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(dynamicConfigPath); os.IsNotExist(err) {
		if err := os.WriteFile(dynamicConfigPath, []byte(defaultConfig), 0755); err != nil {
			return err
		} else {
			return err
		}
	}

	return nil
}
