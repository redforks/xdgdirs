package xdgdirs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestXdgdirs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xdgdirs Suite")
}
