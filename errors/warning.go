package errors

import "strings"

type Warning struct {
	err error
}

func (w *Warning) Error() string {
	return w.err.Error()
}

type Warnings []*Warning

func NewWarnings(errs ...error) Warnings {
	w := Warnings{}
	for _, err := range errs {
		w = append(w, NewWarning(err))
	}
	return w
}

func (errs Warnings) Error() string {
	components := []string{}
	for _, err := range errs {
		components = append(components, err.Error())
	}
	return strings.Join(components, ", ")
}

func (w Warnings) Append(err error) Warnings {
	return append(w, NewWarning(err))
}

func (w Warnings) Any() bool {
	return len(w) > 0
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
	_, isWarning := err.(*Warning)
	_, isWarnings := err.(Warnings)
	return isWarning || isWarnings
}
