package incoming //note also testing unexported methods

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIncoming(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Incoming Suite")
}
