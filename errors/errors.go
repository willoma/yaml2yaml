package errors

import "errors"

var (
	ErrAlreadyDefined                      = errors.New("already defined")
	ErrCannotMergeMappingIntoScalar        = errors.New("cannot merge mapping into scalar")
	ErrCannotMergeMappingIntoSequence      = errors.New("cannot merge mapping into sequence")
	ErrCannotMergeScalarIntoMapping        = errors.New("cannot merge scalar into mapping")
	ErrCannotMergeSequenceIntoMapping      = errors.New("cannot merge sequence into mapping")
	ErrCannotMergeSequenceIntoScalar       = errors.New("cannot merge sequence into scalar")
	ErrEmptySelector                       = errors.New("empty selector")
	ErrInvalidExpression                   = errors.New("invalid expression")
	ErrInvalidOperationWithBinaries        = errors.New("invalid operation with binaries")
	ErrInvalidOperationWithFloats          = errors.New("invalid operation with floats")
	ErrInvalidOperationWithLists           = errors.New("invalid operation with lists")
	ErrInvalidOperationWithMaps            = errors.New("invalid operation with maps")
	ErrInvalidOperationWithStrings         = errors.New("invalid operation with strings")
	ErrInvalidTypeForMathematicalOperation = errors.New("invalid type for mathematical operation")
	ErrNotAMapping                         = errors.New("not a mapping")
	ErrNotAMappingOrSequence               = errors.New("not a mapping or sequence")
	ErrNotASequence                        = errors.New("not a sequence")
	ErrNotDefined                          = errors.New("not defined")
	ErrStreamNotSupported                  = errors.New("kind stream not supported")
	ErrSyntaxError                         = errors.New("syntax error")
	ErrUnexpectedValue                     = errors.New("unexpected value")
	ErrUnkeyedValue                        = errors.New("unkeyed value")
	ErrUnknownKind                         = errors.New("unknown kind")
	ErrUnsupportedExpression               = errors.New("unsupported expression")
	ErrUnsupportedOperator                 = errors.New("unsupported operator")

	ErrNotImplemented = errors.New("not implemented")
)
