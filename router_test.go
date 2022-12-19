package router

import (
	"net/http"
	"testing"
)

func TestMatchingRouteReturnedWhenStaticPathRequested(t *testing.T) {
	type pathCase struct {
		Requested  Route
		Registered []Route
		Wanted     Route
	}

	casesList := []pathCase{
		{
			Requested: NewRoute("/", http.MethodGet, handlerMock),
			Wanted:    NewRoute("/", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/news", http.MethodGet, handlerMock),
			Wanted:    NewRoute("/news", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/comMenTs", http.MethodGet, handlerMock),
			Wanted:    NewRoute("/comments", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/news/334", http.MethodGet, handlerMock),
			Wanted:    NewRoute("/news/:id", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/news/12/comments", http.MethodGet, handlerMock),
			Wanted:    NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
		},
	}

	routesList := []Route{
		NewRoute("/", http.MethodGet, handlerMock),
		NewRoute("/", http.MethodPost, handlerMock),
		NewRoute("/new", http.MethodGet, handlerMock),
		NewRoute("/news", http.MethodGet, handlerMock),
		NewRoute("/news", http.MethodPost, handlerMock),
		NewRoute("/news/:id", http.MethodGet, handlerMock),
		NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
		NewRoute("/c", http.MethodGet, handlerMock),
		NewRoute("/comments", http.MethodGet, handlerMock),
		NewRoute("/comments", http.MethodPost, handlerMock),
		NewRoute("/comments/:id", http.MethodGet, handlerMock),
	}

	router := NewRouter()
	router.AddRoutes(routesList)

	for _, data := range casesList {
		receivedRoute, err := router.GetRoute(data.Requested.GetPath(), data.Requested.GetMethod())
		if nil != err {
			t.Errorf("router did not return any route")
		}
		if receivedRoute.GetPath() != data.Wanted.GetPath() {
			t.Errorf("route with unexpected path received from router")
		}
		if receivedRoute.GetMethod() != data.Wanted.GetMethod() {
			t.Errorf("route with unexpected method received from router")
		}
	}
}

func TestNotFoundWhenNonExistentPathRequested(t *testing.T) {
	wantedRoute := NewRoute("/", http.MethodGet, handlerMock)
	routesList := []Route{
		NewRoute("/v1/news", http.MethodGet, handlerMock),
		NewRoute("/v1/news/:id", http.MethodGet, handlerMock),
		NewRoute("/v1/news/:id/comments", http.MethodGet, handlerMock),
		NewRoute("/v1/comments", http.MethodPost, handlerMock),
	}
	router := NewRouter()
	router.AddRoutes(routesList)
	_, err := router.GetRoute(wantedRoute.GetPath(), wantedRoute.GetMethod())

	if nil == err {
		t.Errorf("got route but non existent path requested")
	}
}

func TestMethodNotAllowedWhenNonExistentMethodRequested(t *testing.T) {
	wantedRoute := NewRoute("/v1/comments", http.MethodGet, handlerMock)
	routesList := []Route{
		NewRoute("/v1/news", http.MethodGet, handlerMock),
		NewRoute("/v1/news/:id", http.MethodGet, handlerMock),
		NewRoute("/v1/comments", http.MethodPost, handlerMock),
	}
	router := NewRouter()
	router.AddRoutes(routesList)

	_, err := router.GetRoute(wantedRoute.GetPath(), wantedRoute.GetMethod())
	if nil == err {
		t.Errorf("got route but not allowed method requested")
	}
}

func TestPanicWhenIntersectingRouteAdded(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("intersecting route was registered")
		}
	}()

	casesList := [][]Route{
		{
			NewRoute("/news", http.MethodGet, handlerMock),
			NewRoute("/news", http.MethodGet, handlerMock),
			NewRoute("/news/:id", http.MethodGet, handlerMock),
			NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
			NewRoute("/news", http.MethodPost, handlerMock),
			NewRoute("/comments", http.MethodPost, handlerMock),
		},
		{
			NewRoute("/news", http.MethodGet, handlerMock),
			NewRoute("/news/:id", http.MethodGet, handlerMock),
			NewRoute("/news/statistics", http.MethodGet, handlerMock),
			NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
			NewRoute("/news", http.MethodPost, handlerMock),
			NewRoute("/comments", http.MethodPost, handlerMock),
		},
		{
			NewRoute("/news", http.MethodGet, handlerMock),
			NewRoute("/news/:id", http.MethodGet, handlerMock),
			NewRoute("/news/:test", http.MethodGet, handlerMock),
			NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
			NewRoute("/news", http.MethodPost, handlerMock),
			NewRoute("/comments", http.MethodPost, handlerMock),
		},
		{
			NewRoute("/news", http.MethodGet, handlerMock),
			NewRoute("/news/:id", http.MethodGet, handlerMock),
			NewRoute("/news/:id/:type", http.MethodGet, handlerMock),
			NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
			NewRoute("/news/example/types", http.MethodGet, handlerMock),
			NewRoute("/news", http.MethodPost, handlerMock),
			NewRoute("/comments", http.MethodPost, handlerMock),
		},
	}
	for _, data := range casesList {
		router := NewRouter()
		router.AddRoutes(data)
	}
}

func handlerMock(http.ResponseWriter, *http.Request) {
}
