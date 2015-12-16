package question_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestQuestion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Question Suite")
}
