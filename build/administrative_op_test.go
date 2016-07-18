package build

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bitbucket.org/atticlab/go-smart-base/xdr"
	"bitbucket.org/atticlab/go-smart-base/keypair"
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
	Describe("tx", func() {
		Context("creating valid tx", func() {
			signer, err := keypair.Random()
			It("created keypair", func() {
				Expect(err).To(BeNil())
			})
			adminOp := AdministrativeOp(OpLongData{OpData: "random_data",})
			It("created admin op", func() {
				Expect(adminOp.Err).To(BeNil())
			})
			tx := Transaction(adminOp, Sequence{1}, SourceAccount{signer.Address()})
			It("created admin op", func() {
				Expect(tx.Err).To(BeNil())
			})
			txE := tx.Sign(signer.Seed())
			rawTxE, err := txE.Base64()
			It("created admin op", func() {
				Expect(err).To(BeNil())
			})
			var newTxE xdr.TransactionEnvelope
			err = xdr.SafeUnmarshalBase64(rawTxE, &newTxE)
			It("created admin op", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
