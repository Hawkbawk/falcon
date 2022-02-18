package darwin

import (
	"fmt"
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
var addResolverCmd string = fmt.Sprintf("echo \"nameserver %v\nport 53\" | sudo tee %v > /dev/null", dnsmasq.LoopbackAddress, dockerResolverPath)
var addLoopbackAddressCmd string = fmt.Sprintf("sudo ifconfig lo0 alias %v", dnsmasq.LoopbackAddress)
var removeResolverCmd string = fmt.Sprintf("sudo rm -f %v", dockerResolverPath)
var removeLoopbackAddressCmd string = fmt.Sprintf("sudo ifconfig lo0 -alias %v", dnsmasq.LoopbackAddress)

// Configure configures the host machine's networking to allow the falcon-proxy to work it's magic.
func Configure() error {
	if err := addDockerResolver(shell.RunCommand); err != nil {
		return fmt.Errorf("Unable to add Docker resolver due to the following error:\n%v", err)
	}

	if err := addLoopbackAddress(shell.RunCommand); err != nil {
		return fmt.Errorf("Unable to add loopback address due to the following error:\n%v", err)
	}

	return nil
}

// Adds the custom *.docker custom resolver.
func addDockerResolver(cmdRunner func(string) error) error {
	logger.LogInfo("Requesting sudo to write to /etc/resolver/docker...")
	if err := cmdRunner(addResolverCmd); err != nil {
		return err
	}
	return nil
}

// Adds the additional loopback address required for inter-container communication to work.
func addLoopbackAddress(cmdRunner func(string) error) error {
	logger.LogInfo("Requesting sudo to add a new loopback address...")
	if err := cmdRunner(addLoopbackAddressCmd); err != nil {
		return err
	}
	return nil
}

// Clean restores the host machine's networking to it's previous state before starting falcon.
func Clean() error {
	if err := removeLoopbackAddress(shell.RunCommand); err != nil {
		return fmt.Errorf("Unable to remove the loopback address due to the following error: \n%v", err)
	}

	if err := removeDockerResolver(shell.RunCommand); err != nil {
		return fmt.Errorf("Unable to remove the custom *.docker resolver due to the following error: \n%v", err)
	}

	return nil
}

// Removes the previously added custom *.docker resolver, if it exists.
func removeDockerResolver(cmdRunner func(string) error) error {
	logger.LogInfo("Requesting sudo to remove /etc/resolver/docker...")
	if err := cmdRunner(removeResolverCmd); err != nil {
		return err
	}
	return nil
}

// Removes the previously added custom loopback address.
func removeLoopbackAddress(cmdRunner func(string) error) error {
	logger.LogInfo("Requesting sudo to remove the added loopback address...")
	err := cmdRunner(removeLoopbackAddressCmd)

	if err != nil {
		if loopbackAlreadyDeletedRegex.Match([]byte(err.Error())) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
