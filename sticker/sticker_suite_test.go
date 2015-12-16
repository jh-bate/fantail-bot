package sticker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSticker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sticker Suite")
}
