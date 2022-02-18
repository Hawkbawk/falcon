package dnsmasq

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDnsmasq(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dnsmasq Suite")
}
