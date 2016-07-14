package build

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AdministrativeOpBuilder Mutators", func() {

	var (
		subject AdministrativeOpBuilder
		mut     interface{}

		validData = "GAXEMCEXBERNSRXOEKD4JAIKVECIXQCENHEBRVSPX2TTYZPMNEDSQCNQ"
		badData   = ""
	)

	JustBeforeEach(func() {
		subject = AdministrativeOpBuilder{}
		subject.Mutate(mut)
	})

	Describe("OpLongData", func() {
		Context("using a valid op data", func() {
			BeforeEach(func() { mut = OpLongData{validData} })

			It("succeeds", func() {
				Expect(subject.Err).NotTo(HaveOccurred())
			})

			It("sets the OpData to the correct value", func() {
				Expect(string(subject.OpData)).To(Equal(validData))
			})
		})

		Context("using an invalid value", func() {
			BeforeEach(func() { mut = OpLongData{badData} })
			It("failed", func() { Expect(subject.Err).To(HaveOccurred()) })
		})
	})
})
