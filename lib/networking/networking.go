package networking

// Configure sets up all networking on the machine for proxying.
func Configure() {
	managerConfigFile := openManagerConfigFile()

	enableDnsmasq(managerConfigFile)
	closeConfigFile(managerConfigFile)
	// Move/backup the resolv file, then create the symlink.
	moveResolvFile()
	letManagerManageResolv()

	createDockerConfFile()
	reloadNetworkManager()
}

// Restore returns all networking on the machine back to it's original state (hopefully)
func Restore() {
	managerConfigFile := openManagerConfigFile()
	disableDnsmasq(managerConfigFile)
	closeConfigFile(managerConfigFile)
	// First remove the symlink, then restore the resolv file.
	stopManagerManagingResolv()
	restoreResolvFile()
	// The docker.conf tells dnsmasq how to resolve requests to *.docker domains.
	deleteDockerConfFile()
	reloadNetworkManager()
}