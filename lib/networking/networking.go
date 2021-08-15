package networking

import (
	"runtime"

	"github.com/Hawkbawk/falcon/lib/logger"
	"github.com/Hawkbawk/falcon/lib/networking/darwin"
	"github.com/Hawkbawk/falcon/lib/networking/linux"
)

// Configure sets up all networking on the machine for proxying.
func Configure() {
	os := runtime.GOOS
	switch os {
	case "linux":
		linux.Configure()
	case "darwin":
		darwin.Configure()
	default:
		logger.LogError("Your current OS of %v is unsupported. We only currently support Ubuntu and macOS.", os)
	}
}

// Restore returns all networking on the machine back to it's original state (hopefully)
func Restore() {
	os := runtime.GOOS
	switch os {
	case "linux":
		linux.Restore()
	case "darwin":
		darwin.Configure()
	default:
		logger.LogError("Your current OS of %v is unsupported. We only currently support Ubuntu and macOS.", os)
	}
}