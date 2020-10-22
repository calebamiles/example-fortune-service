package fortune_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFortune(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fortune Suite")
}
