package router

type NoSuchRouteError struct {
}

func NewNoSuchRouteError() *NoSuchRouteError {
	return &NoSuchRouteError{}
}

func (err NoSuchRouteError) Error() string {
	return "route with requested path does not exist"
}
