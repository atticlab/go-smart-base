package build

import (
	"errors"

	"github.com/atticlab/go-smart-base/amount"
	"github.com/atticlab/go-smart-base/xdr"
)

// Refund groups the creation of a new RefundBuilder with a call to Mutate.
func Refund(muts ...interface{}) (result RefundBuilder) {
	result.Mutate(muts...)
	return
}

// RefundMutator is a interface that wraps the
// MutateRefund operation.  types may implement this interface to
// specify how they modify an xdr.RefundOp object
type RefundMutator interface {
	MutateRefund(interface{}) error
}

// RefundBuilder represents a transaction that is being built.
type RefundBuilder struct {
	O           xdr.Operation
	R           xdr.RefundOp
	Err         error
}

// Mutate applies the provided mutators to this builder's refund or operation.
func (b *RefundBuilder) Mutate(muts ...interface{}) {
	for _, m := range muts {
		var err error
		switch mut := m.(type) {
		case RefundMutator:
			err = mut.MutateRefund(&b.R)
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

// MutateRefund for Asset sets the RefundOp's Asset field
func (m CreditAmount) MutateRefund(o interface{}) (err error) {
	switch o := o.(type) {
	default:
		err = errors.New("Unexpected operation type")
	case *xdr.RefundOp:
		o.Amount, err = amount.Parse(m.Amount)
		if err != nil {
			return
		}

		o.Asset, err = createAlphaNumAsset(m.Code, m.Issuer)
	}
	return
}

// MutateRefund for Original amount
func (m OriginalAmount) MutateRefund(o interface{}) (err error) {
	switch o := o.(type) {
	default:
		err = errors.New("Unexpected operation type")
	case *xdr.RefundOp:
		o.OriginalAmount, err = amount.Parse(m.Amount)
		if err != nil {
			return
		}
	}
	return
}

// MutateRefund for payment ID
func (m PaymentID) MutateRefund(o interface{}) (err error) {
	switch o := o.(type) {
	default:
		err = errors.New("Unexpected operation type")
	case *xdr.RefundOp:
		o.PaymentId = xdr.Int64(m.ID)
		if err != nil {
			return
		}
	}
	return
}

// MutateRefund for Destination sets the RefundOp's Destination field
func (m PaymentSender) MutateRefund(o interface{}) error {
	switch o := o.(type) {
	default:
		return errors.New("Unexpected operation type")
	case *xdr.RefundOp:
		return setAccountId(m.AddressOrSeed, &o.PaymentSource)
	}
	return nil
}
