package router

type noSuchRouteError struct {
}

func newNoSuchRouteError() *noSuchRouteError {
	return &noSuchRouteError{}
}

func (err noSuchRouteError) Error() string {
	return "route with requested path does not exist"
}

type methodNotAllowedError struct {
}

func newMethodNotAllowedError() *methodNotAllowedError {
	return &methodNotAllowedError{}
}

func (err methodNotAllowedError) Error() string {
	return "route with requested method does not exist"
}
