package fortune_test

import (
	"github.com/calebamiles/example-fortune-service/fortune"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FortuneProvider", func() {
	Describe("getting fortunes", func() {
		Context("when fortune is availble from the OS", func() {
			It("returns a fortune from the OS", func() {
				f := fortune.NewProvider([]byte("this isn't an interesting fortune"))

				fortune, err := f.Get()
				Expect(err).ToNot(HaveOccurred(), "expected no error when getting a fortune")
				Expect(fortune).ToNot(BeEmpty())
			})
		})
	})
})
