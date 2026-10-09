package errors

import "errors"

type Joined struct {
	errs []error
}

func NewJoined() *Joined {
	return &Joined{}
}

func (e *Joined) Add(err error) {
	if err == nil || IsEmpty(err) {
		return
	}

	e.errs = append(e.errs, err)
}

func (e *Joined) Compact() error {
	switch len(e.errs) {
	case 0:
		return nil
	case 1:
		return e.errs[0]
	default:
		return errors.Join(e.errs...)
	}
}

func IsEmpty(err error) bool {
	return errors.Is(err, ErrEmpty)
}
