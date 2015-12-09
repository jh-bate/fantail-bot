package testingtsupport_test

import (
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Ω(true).Should(BeTrue())
}
