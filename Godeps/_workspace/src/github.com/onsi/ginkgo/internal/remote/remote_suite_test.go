package remote_test

import (
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestRemote(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Remote Spec Forwarding Suite")
}