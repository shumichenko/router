package router

import (
	"fmt"
	"strings"
)

type Router struct {
	routesList map[string][]Route
}

func NewRouter() *Router {
	router := &Router{
		routesList: map[string][]Route{},
	}

	return router
}

func (r *Router) AddRoutes(routesList []Route) {
	for _, route := range routesList {
		_, err := r.GetRoute(route.GetPath(), route.GetMethod())
		if nil == err {
			panic(fmt.Sprintf(
				"route %s %s intersects with an existing one and cannot be registered", route.GetMethod(), route.GetPath(),
			))
		}

		r.routesList[route.GetMethod()] = append(r.routesList[route.GetMethod()], route)
	}
}

func (r *Router) GetRoute(path string, method string) (Route, error) {
	path = r.formatPath(path)
	if "/" == path {
		return r.findSlashRoute(method)
	}

	requestedParts := r.splitFormattedPath(path)
	if routesToScan, ok := r.routesList[method]; ok {
		route, err := r.findRouteByPathParts(requestedParts, routesToScan)
		if nil == err {
			return route, nil
		}
	}
	for key, group := range r.routesList {
		if key == method {
			continue
		}
		_, err := r.findRouteByPathParts(requestedParts, group)
		if nil == err {
			return Route{}, NewMethodNotAllowedError()
		}
	}

	return Route{}, NewNoSuchRouteError()
}

func (r *Router) findSlashRoute(method string) (Route, error) {
	if routesToScan, ok := r.routesList[method]; ok {
		for _, route := range routesToScan {
			if "/" == route.GetPath() {
				return route, nil
			}
		}
	}

	for key, group := range r.routesList {
		if key == method {
			continue
		}
		for _, route := range group {
			if "/" == route.GetPath() {
				return Route{}, NewMethodNotAllowedError()
			}
		}
	}

	return Route{}, NewNoSuchRouteError()
}

func (r *Router) findRouteByPathParts(requestedParts []string, routesToScan []Route) (Route, error) {
	for _, route := range routesToScan {
		if "/" == route.GetPath() {
			continue
		}
		existingParts := r.splitFormattedPath(route.GetPath())
		if len(existingParts) != len(requestedParts) {
			continue
		}
		if r.isRequestedPathEqual(requestedParts, existingParts) {
			return route, nil
		}
	}

	return Route{}, NewNoSuchRouteError()
}

func (r *Router) isRequestedPathEqual(requestedParts []string, existingParts []string) bool {
	for i, existing := range existingParts {
		isPartStatic := ':' != existing[0]
		if isPartStatic && existing != requestedParts[i] {
			return false
		}
	}

	return true
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
		path = path[0:lastIndex]
	}

	path = strings.ToLower(path)

	return path
}

func (r *Router) splitFormattedPath(path string) []string {
	parts := strings.Split(path, "/")

	return parts[1:]
}
