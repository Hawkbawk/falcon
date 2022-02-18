package darwin

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDarwin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Darwin Suite")
}
