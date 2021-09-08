package darwin

import (
	"fmt"
	"os"
	"regexp"

	"github.com/Hawkbawk/falcon/lib/dnsmasq"
	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/Hawkbawk/falcon/lib/shell"
)

// The path to the resolver file we're going to create to tell macOS to resolve any requests to the
// *.docker to our loopback address.
const dockerResolverPath = "/etc/resolver/docker"

// Allows us to check if, when deleting our new loopback address, the operation failed because
// the address didn't exist in the first place. This way we can make sure falcon down is idempotent
// and doesn't error when run multiple times in a row.
var loopbackAlreadyDeletedRegex regexp.Regexp = *regexp.MustCompile("(SIOCDIFADDR)")

// The command that adds our custom resolver for *.docker domains. By running through the shell, we
// can ask for sudo only when we need it, rather than requiring a user to run falcon with sudo.
var addResolverCmd string = fmt.Sprintf(`echo "nameserver %v" | sudo tee %v > /dev/null`, dnsmasq.LoopbackAddress, dockerResolverPath)

// Configure configures the host machine's networking to allow the falcon-proxy to work it's magic.
func Configure() {
	if err := addDockerResolver(); err != nil {
		logger.LogError("Unable to add Docker resolver due to the following error: \n%v", err)
	}

	if err := addLoopbackAddress(); err != nil {
		logger.LogError("Unable to add loopback address due to the following error: \n%v", err)
	}
}

// Adds the custom *.docker custom resolver.
func addDockerResolver() error {
	logger.LogInfo("Requesting sudo to write to /etc/resolver/docker...")
	if err := shell.RunCommand(addResolverCmd); err != nil {
		return err
	}
	return nil
}

// Adds the additional loopback address required for inter-container communication to work.
func addLoopbackAddress() error {
	logger.LogInfo("Requesting sudo to add a new loopback address...")
	if err := shell.RunCommand(fmt.Sprintf("sudo ifconfig lo0 alias %v", dnsmasq.LoopbackAddress)); err != nil {
		return err
	}
	return nil
}

// Clean restores the host machine's networking to it's previous state before starting falcon.
func Clean() {
	if err := removeLoopbackAddress(); err != nil {
		logger.LogError("Unable to remove the loopback address due to the following error: \n%v", err)
	}

	if err := removeDockerResolver(); err != nil {
		logger.LogError("Unable to remove the custom *.docker resolver due to the following error: \n%v", err)
	}
}

// Removes the previously added custom *.docker resolver, if it exists.
func removeDockerResolver() error {
	_, err := os.Stat(dockerResolverPath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	logger.LogInfo("Requesting sudo to remove /etc/resolver/docker...")
	if err := shell.RunCommand(fmt.Sprintf("sudo rm -f %v", dockerResolverPath)); err != nil {
		return err
	}
	return nil
}

// Removes the previously added custom loopback address.
func removeLoopbackAddress() error {
	logger.LogInfo("Requesting sudo to remove the added loopback address...")
	err := shell.RunCommand(fmt.Sprintf("sudo ifconfig lo0 -alias %v", dnsmasq.LoopbackAddress))

	if err != nil {
		if loopbackAlreadyDeletedRegex.Match([]byte(err.Error())) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
