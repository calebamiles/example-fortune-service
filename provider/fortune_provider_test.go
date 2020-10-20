package provider_test

import (
	"os"

	"github.com/calebamiles/example-fortune-service/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FortuneProvider", func() {
	Describe("getting fortunes", func() {
		Context("when fortune is availble from the OS", func() {
			It("returns a fortune from the OS", func() {
				f := provider.NewFortune([]byte("this isn't an interesting fortune"))

				fortune, err := f.Get()
				Expect(err).ToNot(HaveOccurred(), "expected no error when getting a fortune")
				Expect(fortune).ToNot(BeEmpty())
			})
		})

		Context("when fortune is not available from the OS", func() {
			It("returns the default fortune", func() {
				oldPath := os.Getenv("PATH")
				defer os.Setenv("PATH", oldPath)

				// Clear the PATH
				os.Setenv("PATH", "")

				defaultFortune := []byte("this isn't an interesting fortune")
				f := provider.NewFortune(defaultFortune)

				fortune, err := f.Get()
				Expect(err).ToNot(HaveOccurred(), "expected no error when getting a fortune")
				Expect(fortune).To(Equal(fortune))
			})
		})
	})
})
