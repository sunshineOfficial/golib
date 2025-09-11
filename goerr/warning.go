package goerr

type Warning struct {
	origin  error
	message string
}

func NewWarning(message string, err error) Warning {
	return Warning{
		origin:  err,
		message: message,
	}
}

func (n Warning) Error() string {
	if n.origin == nil {
		return n.message
	}

	return n.origin.Error()
}

func (n Warning) Unwrap() error {
	return n.origin
}
