package errors

import "strings"

type MultiError struct {
	errs []error
}

func Multi() *MultiError {
	return &MultiError{}
}

func (m *MultiError) IsNil() bool {
	return len(m.errs) == 0
}

func (m *MultiError) Add(err error) {
	if err == nil {
		return
	}

	//nolint:errorlint // We need to check precisely this err, not its wrapped errors recursively.
	if mErr, ok := err.(*MultiError); ok {
		m.errs = append(m.errs, mErr.errs...)

		return
	}

	m.errs = append(m.errs, err)
}

func (m *MultiError) Error() string {
	switch len(m.errs) {
	case 0:
		return ""
	case 1:
		return m.errs[0].Error()
	default:
		var str strings.Builder

		str.WriteString(m.errs[0].Error())

		for _, err := range m.errs[1:] {
			str.WriteString("\n")
			str.WriteString(err.Error())
		}

		return str.String()
	}
}

func (m *MultiError) Unwrap() []error {
	return m.errs
}
