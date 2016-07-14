package build

import (
	"errors"

	"bitbucket.org/atticlab/go-smart-base/xdr"
)

// AdministrativeOp groups the creation of a new AdministrativeOpBuilder with a call to Mutate.
func AdministrativeOp(muts ...interface{}) (result AdministrativeOpBuilder) {
	result.Mutate(muts...)
	return
}

// AdministrativeOpMutator is a interface that wraps the
// MutateAdministrativeOp operation.  types may implement this interface to
// specify how they modify an xdr.AdministrativeOpBuilder object
type AdministrativeOpMutator interface {
	MutateAdministrativeOp(*AdministrativeOpBuilder) error
}

// AdministrativeOpBuilder represents a transaction that is being built.
type AdministrativeOpBuilder struct {
	O      xdr.Operation
	OpData xdr.LongString
	Err    error
}

// Mutate applies the provided mutators to this builder or operation.
func (b *AdministrativeOpBuilder) Mutate(muts ...interface{}) {
	for _, m := range muts {
		var err error
		switch mut := m.(type) {
		case AdministrativeOpMutator:
			err = mut.MutateAdministrativeOp(b)
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

// MutateAdministrativeOp for OpData sets the AdministrativeOpBuilder's OpData field
func (m OpLongData) MutateAdministrativeOp(o *AdministrativeOpBuilder) error {
	if m.OpData == "" {
		return errors.New("OpData can't be empty")
	}
	o.OpData = xdr.LongString(m.OpData)
	return nil
}
