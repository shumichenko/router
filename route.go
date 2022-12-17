package router

import (
	"net/http"
	"strings"
)

type Route struct {
	Path    string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

func NewRoute(path string, method string, handler func(http.ResponseWriter, *http.Request)) Route {
	if len(path) < 1 {
		panic("Path cannot be empty")
	}
	if path[0] != '/' {
		panic("Path should start with /")
	}
	if len(method) < 1 {
		panic("Method cannot be empty")
	}

	lastIndex := len(path) - 1
	if lastIndex > 0 && '/' == path[lastIndex] {
		path = path[0:lastIndex]
	}
	path = strings.ToLower(path)

	return Route{
		Path:    path,
		Method:  method,
		Handler: handler,
	}
}
