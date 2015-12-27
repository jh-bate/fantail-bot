package sticker_test

import (
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestSticker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sticker Suite")
}
