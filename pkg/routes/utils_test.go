package routes

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	BeforeEach(func() {

	})

	Describe("EnsurePrivatekey", func() {
		It("should error if no account is provided", func() {
			p, err := EnsurePrivateKey("")
			Expect(p).To(BeNil())
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("account not passed"))
		})
		It("should not error if account is provided", func() {
			p, err := EnsurePrivateKey("abcd")
			Expect(p).NotTo(BeNil())
			Expect(err).To(BeNil())
		})
	})
})
