package bot //note also testing unexported methods

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bot Suite")
}
