package errors

type ValidationError struct {
	E error
}

func (err *ValidationError) Error() string {
	return err.E.Error()
}
