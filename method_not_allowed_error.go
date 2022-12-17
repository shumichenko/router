package router

type MethodNotAllowedError struct {
}

func NewMethodNotAllowedError() *MethodNotAllowedError {
	return &MethodNotAllowedError{}
}

func (err MethodNotAllowedError) Error() string {
	return "route with requested method does not exist"
}
