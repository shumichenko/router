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
			Registered: []Route{
				NewRoute("/", http.MethodGet, handlerMock),
				NewRoute("/new", http.MethodGet, handlerMock),
				NewRoute("/news/:news", http.MethodGet, handlerMock),
				NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
				NewRoute("/news", http.MethodGet, handlerMock),
				NewRoute("/news", http.MethodPost, handlerMock),
				NewRoute("/", http.MethodPost, handlerMock),
			},
			Wanted: NewRoute("/", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/news", http.MethodGet, handlerMock),
			Registered: []Route{
				NewRoute("/", http.MethodGet, handlerMock),
				NewRoute("/new", http.MethodGet, handlerMock),
				NewRoute("/news/:news", http.MethodGet, handlerMock),
				NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
				NewRoute("/news", http.MethodGet, handlerMock),
				NewRoute("/news", http.MethodPost, handlerMock),
				NewRoute("/comments", http.MethodPost, handlerMock),
			},
			Wanted: NewRoute("/news", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/comMenTs", http.MethodGet, handlerMock),
			Registered: []Route{
				NewRoute("/", http.MethodGet, handlerMock),
				NewRoute("/c", http.MethodGet, handlerMock),
				NewRoute("/comments", http.MethodGet, handlerMock),
				NewRoute("/comments/:id", http.MethodGet, handlerMock),
				NewRoute("/news/:id/comments", http.MethodGet, handlerMock),
				NewRoute("/news", http.MethodPost, handlerMock),
				NewRoute("/comments", http.MethodPost, handlerMock),
			},
			Wanted: NewRoute("/comments", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/v1/news/334", http.MethodGet, handlerMock),
			Registered: []Route{
				NewRoute("/v1/news", http.MethodGet, handlerMock),
				NewRoute("/v1/news/:id", http.MethodGet, handlerMock),
				NewRoute("/v1/news/:id/comments", http.MethodGet, handlerMock),
				NewRoute("/v1/comments", http.MethodPost, handlerMock),
			},
			Wanted: NewRoute("/v1/news/:id", http.MethodGet, handlerMock),
		},
		{
			Requested: NewRoute("/v1/news/12/comments", http.MethodGet, handlerMock),
			Registered: []Route{
				NewRoute("/v1/news", http.MethodGet, handlerMock),
				NewRoute("/v1/news/:id", http.MethodGet, handlerMock),
				NewRoute("/v1/news/:id/comments", http.MethodGet, handlerMock),
				NewRoute("/v1/comments", http.MethodPost, handlerMock),
			},
			Wanted: NewRoute("/v1/news/:id/comments", http.MethodGet, handlerMock),
		},
	}

	for _, data := range casesList {
		router := NewRouter()
		router.AddRoutes(data.Registered)

		receivedRoute, err := router.GetRoute(data.Requested.Path, data.Requested.Method)
		if nil != err {
			t.Errorf("Router did not return any route")
		}
		if receivedRoute.Path != data.Wanted.Path {
			t.Errorf("Route with unexpected path received from router")
		}
		if receivedRoute.Method != data.Wanted.Method {
			t.Errorf("Route with unexpected method received from router")
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
	_, err := router.GetRoute(wantedRoute.Path, wantedRoute.Method)

	if nil == err {
		t.Errorf("Got route but non existent path requested")
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

	_, err := router.GetRoute(wantedRoute.Path, wantedRoute.Method)
	if nil == err {
		t.Errorf("Got route but not allowed method requested")
	}
}

func TestPanicWhenIntersectingRouteAdded(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("Intersecting route was registered")
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
	}
	for _, data := range casesList {
		router := NewRouter()
		router.AddRoutes(data)
	}
}

func handlerMock(http.ResponseWriter, *http.Request) {
}
