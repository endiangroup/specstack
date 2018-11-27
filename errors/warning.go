package errors

type Warnings Errors

type Warning struct {
	err error
}

func (w *Warning) Error() string {
	return w.err.Error()
}

func NewWarning(err error) *Warning {
	return &Warning{err}
}

func WarningOrNil(err error) error {
	if err == nil {
		return nil
	}

	return NewWarning(err)
}

func IsWarning(err error) bool {
	_, ok := err.(*Warning)
	return ok
}
