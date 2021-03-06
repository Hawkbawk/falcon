// The networkmanager package contains all of the functions necessary to set up a Linux machine
// to enable NetworkManager dnsmasq and let NetworkManager control resolv.conf
package linux

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/Hawkbawk/falcon/lib/dnsmasq"
	"github.com/Hawkbawk/falcon/lib/files"
	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/Hawkbawk/falcon/lib/shell"
)

const managerConfigFilePath = "/etc/NetworkManager/NetworkManager.conf"
const managerResolvFilePath = "/var/run/NetworkManager/resolv.conf"
const dockerConfFilePath = "/etc/NetworkManager/dnsmasq.d/docker.conf"
const dnsmasqLine = "dns=dnsmasq\n"

// This is the secret sauce that allows falcon to work better than dory. This new
// loopback address ensures that intercontainer communication works properly, as requests
// to *.docker will resolve to this address, rather than 127.0.0.1. This ensures that all
// requests from one container to another actually go through the proxy and don't get stuck
// inside the container itself.
const loopbackAddress = "192.168.40.1"
const netmask = "32"

// This line will need to be updated, as we no longer set a static IP for the proxy, which means
// we need to dynamically determine the IP address.
var dockerConfLine string = fmt.Sprint("address=/docker/", loopbackAddress)
var mainSectionRegex *regexp.Regexp = regexp.MustCompile(`\[main\]`)
var dnsmasqEnabledRegex *regexp.Regexp = regexp.MustCompile(`^dns=dnsmasq$`)

// OpenManagerConfigFile gets a file handler for the NetworkManager config file.
func openManagerConfigFile() *os.File {
	return files.OpenFile(managerConfigFilePath)
}

// closeConfigFile closes the file handler for the configFile. It's a useful helper method
// for deferring cleanup to. If it encounters any errors, it ends the program.
func closeConfigFile(configFile *os.File) {
	err := configFile.Close()

	if err != nil {
		logger.LogError("Unable to close config file. This is likely due to a bug.")
	}
}

// EnableDnsmasq writes the necessary line to the NetworkManager.conf file to enable dnsmasq
// for the machine.
func enableDnsmasq(configFile *os.File) {
	if dnsMasqEnabled(configFile) {
		return
	}

	previousContents := files.ReadFile(configFile)
	newConfigFile := make([]byte, files.FileSize(configFile)+int64(len(dnsmasqLine)))

	indices := mainSectionRegex.FindIndex(previousContents)

	if indices == nil {
		logger.LogError("You don't have a main section in your NetworkManager.conf! You should probably add one.")
	}

	beforeMain := previousContents[:indices[0]]
	mainLine := previousContents[indices[0]:indices[1]]
	afterMain := previousContents[indices[1]:]

	// Copy everything into this new config file, including the line that adds dnsmasq.
	copy(newConfigFile, beforeMain)
	copy(newConfigFile, mainLine)
	copy(newConfigFile, []byte(dnsmasqLine))
	copy(newConfigFile, afterMain)

	files.OverwriteFile(configFile, newConfigFile)
}

// DisableDnsmasq removes the line necessary in NetworkManager.conf to enable dnsmasq for
// the system.
func disableDnsmasq(configFile *os.File) {
	if !dnsMasqEnabled(configFile) {
		return
	}

	previousContents := files.ReadFile(configFile)
	restoredConfigFile := make([]byte, files.FileSize(configFile)-int64(len(dnsmasqLine)))
	// Indices cannot be nil, as we have already checked to see if dnsmasq was disabled and
	// both methods use the same regex.
	indices := dnsmasqEnabledRegex.FindIndex(previousContents)

	beforeDnsmasq := previousContents[:indices[0]]
	afterDnsmasq := previousContents[indices[1]:]

	copy(restoredConfigFile, beforeDnsmasq)
	copy(restoredConfigFile, afterDnsmasq)

	files.OverwriteFile(configFile, restoredConfigFile)
}

// dnsmasqEnabled determines whether dnsmasq through NetworkManager has already been enabled.
func dnsMasqEnabled(configFile *os.File) bool {
	fileContents := files.ReadFile(configFile)

	return dnsmasqEnabledRegex.Match(fileContents)
}

func letManagerManageResolv() {
	files.Symlink(managerResolvFilePath, resolvFilePath)
}

func stopManagerManagingResolv() {
	files.DeleteFile(resolvFilePath)
}

func createDockerConfFile() {
	dockerConfFile := files.CreateFile(dockerConfFilePath)

	io.WriteString(dockerConfFile, dockerConfLine)

	if err := dockerConfFile.Close(); err != nil {
		logger.LogError("Unable to create dnsmasq docker.conf file. Error: ", err.Error())
	}
}

func deleteDockerConfFile() {
	files.DeleteFile(dockerConfFilePath)
}

func addLoopbackAddress() {
	if err := shell.RunCommand(fmt.Sprintf("sudo ip addr add %v/%v dev lo", dnsmasq.LoopbackAddress, netmask)); err != nil {
		logger.LogError("Unable to add the loopback address required by falcon. ERROR: ", err)
	}
}

func removeLoopbackAddress() {
	if err := shell.RunCommand(fmt.Sprintf("sudo ip addr del %v/%v dev lo", dnsmasq.LoopbackAddress, netmask)); err != nil {
		logger.LogError("Unable to remove the loopback address added by falcon. ERROR: ", err)

	}
}

func reloadNetworkManager() {
	if err := shell.RunCommand("sudo systemctl reload NetworkManager"); err != nil {
		logger.LogError("Unable to restart NetworkManager. Error: ", err.Error())
	}
}
