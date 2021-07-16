package resolv

import (
	"fmt"
	"os"

	"github.com/hawkbawk/prox/lib/logger"
)

const resolvFilePath = "/etc/resolv.conf"
const backupFilePath = "~/go/github.com/falcon/backups"


func readFile(filePath string, fileName string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("couldn't open file: %v", fileName)
	}

	return string(data), nil
}

func moveResolvFile() {
	err := os.Rename(resolvFilePath, backupFilePath)

	if err != nil {
		logger.LogError("Unable to backup your resolv file.")
	}
}

