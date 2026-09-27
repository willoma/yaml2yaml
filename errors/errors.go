package errors

import "strings"

type multiErrors []error

func (e multiErrors) Error() string {
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
	for _, e := range errs {
		if e != nil {
			result = append(result, e)
		}
	}

	switch len(result) {
	case 0:
		return nil
	case 1:
		return result[0]
	default:
		return multiErrors(result)
	}
}
