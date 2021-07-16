package networking

import (
	"os"

	"github.com/hawkbawk/falcon/lib/logger"
)

const resolvFilePath = "/etc/resolv.conf"
const backupFilePath = "~/go/github.com/falcon/backups"

// Moves the current resolv.conf file into the backup directory located at
// ~/go/github.com/falcon/backups. In the future, a more universal backup location will likely
// be added.
func moveResolvFile() {
	err := os.Rename(resolvFilePath, backupFilePath)

	if err != nil {
		logger.LogError("Unable to backup your resolv file.")
	}
}

// Moves the backed up resolv.conf file back to it's usual spot at /etc/resolv.conf.
func restoreResolvFile() {
	err := os.Rename(backupFilePath, resolvFilePath)

	if err != nil {
		logger.LogError("Unable to restore your backed up resolve file. It may have been deleted.")
	}
}
