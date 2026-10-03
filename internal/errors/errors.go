package errors

import "errors"

var (
	Join = errors.Join

	ErrEmptyNode                    = errors.New("empty node")
	ErrEmptySelector                = errors.New("empty selector")
	ErrInvalidExpression            = errors.New("invalid expression")
	ErrInvalidOperationWithBinaries = errors.New("invalid operation with binary data")
	ErrInvalidOperationWithFloats   = errors.New("invalid operation with floats")
	ErrInvalidOperationWithLists    = errors.New("invalid operation with lists")
	ErrInvalidOperationWithMaps     = errors.New("invalid operation with maps")
	ErrInvalidOperationWithStrings  = errors.New("invalid operation with strings")
	ErrNotAMapping                  = errors.New("not a mapping")
	ErrNotAMappingEntry             = errors.New("not a mapping entry")
	ErrNotAMappingOrSequence        = errors.New("not a mapping or sequence")
	ErrNotASequence                 = errors.New("not a sequence")
	ErrNotDefined                   = errors.New("not defined")
	ErrSyntaxError                  = errors.New("syntax error")
	ErrUnkeyedValue                 = errors.New("unkeyed value")
	ErrUnsupportedExpression        = errors.New("unsupported expression")
	ErrUnsupportedOperator          = errors.New("unsupported operator")

// ErrNotImplemented = errors.New("not implemented")
)
