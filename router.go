package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type PathParams []PathParameter

func (params PathParams) GetByName(name string) string {
	for _, parameter := range params {
		if parameter.Key == name {
			return parameter.Value
		}
	}

	return ""
}

func GetParamsFromContext(ctx context.Context) PathParams {
	p, _ := ctx.Value(ParamsKey).(PathParams)

	return p
}

type PathParameter struct {
	Key   string
	Value string
}

type paramsKey struct{}

var ParamsKey = paramsKey{}

type Router struct {
	routesList       map[string][]Route
	NotFound         http.Handler
	MethodNotAllowed http.Handler
	PanicHandler     func(http.ResponseWriter, *http.Request, interface{})
}

func NewRouter() *Router {
	router := &Router{
		routesList: map[string][]Route{},
	}

	return router
}

func (r *Router) AddRoutes(routesList []Route) {
	for _, route := range routesList {
		foundRoute, _, err := r.getRoute(route.GetPath(), route.GetMethod())
		if nil == err {
			panic(fmt.Sprintf(
				"route %s %s intersects with an existing one %s %s and cannot be registered",
				route.GetMethod(),
				route.GetPath(),
				foundRoute.GetMethod(),
				foundRoute.GetPath(),
			))
		}
		r.routesList[route.GetMethod()] = append(r.routesList[route.GetMethod()], route)
	}
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if nil != r.PanicHandler {
		defer r.recoverPanic(writer, request)
	}

	route, params, err := r.getRoute(request.URL.RequestURI(), request.Method)
	switch err.(type) {
	case *noSuchRouteError:
		r.handleNotFound(writer, request)

		return
	case *methodNotAllowedError:
		r.handleMethodNotAllowed(writer, request)

		return
	}

	if len(params) > 0 {
		ctx := request.Context()
		ctx = context.WithValue(ctx, ParamsKey, params)
		request = request.WithContext(ctx)
	}
	handler := route.GetHandler()
	handler(writer, request)
}

func (r *Router) getRoute(path string, method string) (*Route, PathParams, error) {
	path = r.formatPath(path)

	requestedParts := r.splitFormattedPath(path)
	if routesToScan, ok := r.routesList[method]; ok {
		route, params := r.findRoute(routesToScan, requestedParts)
		if nil != route {
			return route, params, nil
		}
	}
	for key, group := range r.routesList {
		if key == method {
			continue
		}
		route, params := r.findRoute(group, requestedParts)
		if nil != route {
			return nil, params, newMethodNotAllowedError()
		}
	}

	return nil, PathParams{}, newNoSuchRouteError()
}

func (r *Router) findRoute(routesToScan []Route, requestedParts []string) (*Route, PathParams) {
	if "" == requestedParts[0] {
		route := r.findSlashRoute(routesToScan)

		return route, PathParams{}
	}

	return r.findRouteByPathParts(routesToScan, requestedParts)
}

func (r *Router) findSlashRoute(routesToScan []Route) *Route {
	for _, route := range routesToScan {
		if "/" == route.GetPath() {
			return &route
		}
	}

	return nil
}

func (r *Router) findRouteByPathParts(routesToScan []Route, requestedParts []string) (*Route, PathParams) {
	for _, route := range routesToScan {
		if "/" == route.GetPath() {
			continue
		}
		existingParts := r.splitFormattedPath(route.GetPath())
		if len(existingParts) != len(requestedParts) {
			continue
		}
		if ok, params := r.isRequestedPathEqual(requestedParts, existingParts); ok {
			return &route, params
		}
	}

	return nil, PathParams{}
}

func (r *Router) isRequestedPathEqual(requestedParts []string, existingParts []string) (bool, PathParams) {
	var params PathParams
	for i, existing := range existingParts {
		isPartStatic := ':' != existing[0]
		if !isPartStatic {
			params = append(params, PathParameter{existing[1:], requestedParts[i]})

			continue
		} else if isPartStatic && existing != requestedParts[i] {
			return false, PathParams{}
		}
	}

	return true, params
}

func (r *Router) formatPath(path string) string {
	if len(path) < 1 {
		return "/"
	}

	if '/' != path[0] {
		path = "/" + path
	}
	lastIndex := len(path) - 1
	if lastIndex > 0 && '/' == path[lastIndex] {
		path = path[:lastIndex]
	}

	return strings.ToLower(path)
}

func (r *Router) splitFormattedPath(path string) []string {
	parts := strings.Split(path, "?")
	parts = strings.Split(parts[0], "/")

	return parts[1:]
}

func (r *Router) recoverPanic(writer http.ResponseWriter, request *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(writer, request, rcv)
	}
}

func (r *Router) handleNotFound(writer http.ResponseWriter, request *http.Request) {
	if nil == r.NotFound {
		http.NotFound(writer, request)
	} else {
		r.NotFound.ServeHTTP(writer, request)
	}
}

func (r *Router) handleMethodNotAllowed(writer http.ResponseWriter, request *http.Request) {
	if nil == r.MethodNotAllowed {
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	} else {
		r.MethodNotAllowed.ServeHTTP(writer, request)
	}
}
