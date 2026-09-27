package errors

import "strings"

type multi []error

func (e multi) Error() string {
	var buf strings.Builder

	buf.WriteString(e[0].Error())

	for _, err := range e[1:] {
		buf.WriteString(", ")
		buf.WriteString(err.Error())
	}

	return buf.String()
}

func Join(errs ...error) error {
	result := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			if mErr, ok := err.(multi); ok {
				result = append(result, mErr...)
			} else {
				result = append(result, err)
			}
		}
	}

	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0]
	default:
		return multi(result)
	}
}
