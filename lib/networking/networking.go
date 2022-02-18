package networking

import (
	"fmt"
	"runtime"

	"github.com/Hawkbawk/falcon/lib/networking/darwin"
)

// Configure sets up all networking on the machine for proxying.
func Configure() error {
	os := runtime.GOOS
	switch os {
	// Linux support is disabled currently, but will be fixed later
	// case "linux":
	// 	return linux.Configure()
	case "darwin":
		return darwin.Configure()
	default:
		return fmt.Errorf("we only support macOS currently")
	}
}

// Clean returns all networking on the machine back to it's original state (hopefully)
func Clean() error {
	os := runtime.GOOS
	switch os {
	// Linux support is disabled currently, but will be fixed later
	// case "linux":
	// 	return linux.Restore()
	case "darwin":
		return darwin.Clean()
	default:
		return fmt.Errorf("we only support macOS currently")
	}
}
