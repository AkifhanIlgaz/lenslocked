package errors

/*
Public wraps the original error with a new error that has a `Public() string` method
that will return a message that is acceptable to display to the public.
This error can also be unwrapped using the traditional `errors` package approach.
*/
func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

func (pErr publicError) Error() string {
	return pErr.err.Error()
}

func (pErr publicError) Public() string {
	return pErr.msg
}

func (pErr publicError) Unwrap() error {
	return pErr.err
}
