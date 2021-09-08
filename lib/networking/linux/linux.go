package linux

func Configure() {
	managerConfigFile := openManagerConfigFile()

	addLoopbackAddress()
	enableDnsmasq(managerConfigFile)
	// We can't defer this, as we have to save the change we make before we reload NetworkManager.
	closeConfigFile(managerConfigFile)

	// Move/backup the resolv file, then create the symlink.
	moveResolvFile()
	letManagerManageResolv()

	createDockerConfFile()
	reloadNetworkManager()
}

func Restore() {
	managerConfigFile := openManagerConfigFile()

	removeLoopbackAddress()
	disableDnsmasq(managerConfigFile)
	// We can't defer this, as we have to save the change we make before we reload NetworkManager.
	closeConfigFile(managerConfigFile)
	// First remove the symlink, then restore the resolv file.
	stopManagerManagingResolv()
	restoreResolvFile()
	deleteDockerConfFile()
	reloadNetworkManager()
}
