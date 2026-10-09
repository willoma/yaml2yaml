package errors

import "errors"

var (
	ErrEmpty                         = errors.New("empty")
	ErrInvalidExpression             = errors.New("invalid expression")
	ErrInvalidOperationWithBinaries  = errors.New("invalid operation with binary data")
	ErrInvalidOperationWithFloats    = errors.New("invalid operation with floats")
	ErrInvalidOperationWithIntegers  = errors.New("invalid operation with integers")
	ErrInvalidOperationWithSequences = errors.New("invalid operation with sequences")
	ErrInvalidOperationWithMaps      = errors.New("invalid operation with maps")
	ErrInvalidOperationWithStrings   = errors.New("invalid operation with strings")
	ErrNotABinary                    = errors.New("not a binary")
	ErrNotAMapping                   = errors.New("not a mapping")
	ErrNotAMappingEntry              = errors.New("not a mapping entry")
	ErrNotAMappingOrSequence         = errors.New("not a mapping or sequence")
	ErrNotANumber                    = errors.New("not a number")
	ErrNotAPath                      = errors.New("not a path")
	ErrNotASequence                  = errors.New("not a sequence")
	ErrNotDefined                    = errors.New("not defined")
	ErrOutOfBounds                   = errors.New("out of bounds")
	ErrSyntaxError                   = errors.New("syntax error")
	ErrUnkeyedValue                  = errors.New("unkeyed value")
	ErrUnsupportedExpression         = errors.New("unsupported expression")
	ErrUnsupportedOperator           = errors.New("unsupported operator")
)
