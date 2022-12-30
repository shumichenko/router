package router

import (
	"net/http"
	"strings"
)

type Route struct {
	path    string
	method  string
	handler func(http.ResponseWriter, *http.Request)
}

func NewRoute(path string, method string, handler func(http.ResponseWriter, *http.Request)) Route {
	if len(path) < 1 {
		panic("path cannot be empty")
	}
	if path[0] != '/' {
		panic("path should start with /")
	}
	if len(method) < 1 {
		panic("method cannot be empty")
	}
	if strings.Contains("?", path) {
		panic("path should not contain ? character ")
	}
	lastIndex := len(path) - 1
	if lastIndex > 0 && '/' == path[lastIndex] {
		panic("path should not contain trailing /")
	}

	path = strings.ToLower(path)

	return Route{
		path:    path,
		method:  method,
		handler: handler,
	}
}

func (r Route) GetPath() string {
	return r.path
}

func (r Route) GetMethod() string {
	return r.method
}

func (r Route) GetHandler() func(http.ResponseWriter, *http.Request) {
	return r.handler
}
