package build

import (
	"errors"

	"bitbucket.org/atticlab/go-smart-base/amount"
	"bitbucket.org/atticlab/go-smart-base/xdr"
)

// PaymentReversal groups the creation of a new PaymentReversalBuilder with a call to Mutate.
func PaymentReversal(muts ...interface{}) (result PaymentReversalBuilder) {
	result.Mutate(muts...)
	return
}

// PaymentReversalMutator is a interface that wraps the
// MutatePaymentReversal operation.  types may implement this interface to
// specify how they modify an xdr.PaymentReversalOp object
type PaymentReversalMutator interface {
	MutatePaymentReversal(interface{}) error
}

// PaymentReversalBuilder represents a transaction that is being built.
type PaymentReversalBuilder struct {
	O           xdr.Operation
	P           xdr.PaymentReversalOp
	Err         error
}

// Mutate applies the provided mutators to this builder's payment reversal or operation.
func (b *PaymentReversalBuilder) Mutate(muts ...interface{}) {
	for _, m := range muts {
		var err error
		switch mut := m.(type) {
		case PaymentReversalMutator:
			err = mut.MutatePaymentReversal(&b.P)
		case OperationMutator:
			err = mut.MutateOperation(&b.O)
		default:
			err = errors.New("Mutator type not allowed")
		}

		if err != nil {
			b.Err = err
			return
		}
	}
}

// MutatePaymentReversal for Asset sets the PaymentReversalOp's Asset field
func (m CreditAmount) MutatePaymentReversal(o interface{}) (err error) {
	switch o := o.(type) {
	default:
		err = errors.New("Unexpected operation type")
	case *xdr.PaymentReversalOp:
		o.Amount, err = amount.Parse(m.Amount)
		if err != nil {
			return
		}

		o.Asset, err = createAlphaNumAsset(m.Code, m.Issuer)
	}
	return
}

// MutatePaymentReversal for Commission amount
func (m CommissionAmount) MutatePaymentReversal(o interface{}) (err error) {
	switch o := o.(type) {
	default:
		err = errors.New("Unexpected operation type")
	case *xdr.PaymentReversalOp:
		o.CommissionAmount, err = amount.Parse(m.Amount)
		if err != nil {
			return
		}
	}
	return
}

// MutatePaymentReversal for payment ID
func (m PaymentID) MutatePaymentReversal(o interface{}) (err error) {
	switch o := o.(type) {
	default:
		err = errors.New("Unexpected operation type")
	case *xdr.PaymentReversalOp:
		o.PaymentId = xdr.Int64(m.ID)
		if err != nil {
			return
		}
	}
	return
}

// MutatePaymentReversal for Destination sets the PaymentReversalOp's Destination field
func (m PaymentSender) MutatePaymentReversal(o interface{}) error {
	switch o := o.(type) {
	default:
		return errors.New("Unexpected operation type")
	case *xdr.PaymentReversalOp:
		return setAccountId(m.AddressOrSeed, &o.PaymentSource)
	}
	return nil
}
