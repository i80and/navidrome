package utils_test

import (
	"github.com/deluan/navidrome/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AtomicBool", func() {
	var b utils.AtomicBool

	BeforeEach(func() {
		b = utils.AtomicBool{}
	})

	It("initializes with value = false", func() {
		Expect(b.Get()).To(BeFalse())
	})

	It("sets value", func() {
		b.Set(true)
		Expect(b.Get()).To(BeTrue())

		b.Set(false)
		Expect(b.Get()).To(BeFalse())
	})
})
