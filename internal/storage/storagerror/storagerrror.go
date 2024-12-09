package storagerror

type RetriableError struct {
	err error
}

func (re *RetriableError) Error() string {
	return re.err.Error()
}
func (re *RetriableError) Unwrap() error {
	return re.err
}

func NewRetriableError(err error) *RetriableError {
	return &RetriableError{
		err: err,
	}
}
