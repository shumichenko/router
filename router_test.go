package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMatchingRouteReturnedWhenExistingPathRequested(t *testing.T) {
	type pathCase struct {
		Requested Route
		Wanted    Route
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
			Requested: NewRoute("/news?example=true&something=else", http.MethodGet, handlerMock),
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

	createRoute := func(path string, method string) Route {
		return NewRoute(path, method, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte(
				fmt.Sprintf("%s %s", method, path),
			))
		})
	}

	routesList := []Route{
		createRoute("/", http.MethodGet),
		createRoute("/", http.MethodPost),
		createRoute("/new", http.MethodGet),
		createRoute("/news", http.MethodGet),
		createRoute("/news", http.MethodPost),
		createRoute("/news/:id", http.MethodGet),
		createRoute("/news/:id/comments", http.MethodGet),
		createRoute("/c", http.MethodGet),
		createRoute("/comments", http.MethodGet),
		createRoute("/comments", http.MethodPost),
		createRoute("/comments/:id", http.MethodGet),
	}

	router := NewRouter()
	router.AddRoutes(routesList)

	for _, data := range casesList {
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(data.Requested.GetMethod(), data.Requested.GetPath(), nil)
		router.ServeHTTP(recorder, request)
		content, _ := io.ReadAll(recorder.Body)

		result := strings.Split(string(content), " ")
		if len(result) < 1 {
			t.Errorf("router did not return any route")
		} else if data.Wanted.GetMethod() != result[0] {
			t.Errorf("route with unexpected method received from router")
		} else if data.Wanted.GetPath() != result[1] {
			t.Errorf("route with unexpected path received from router")
		}
	}
}

func TestValidPathParamsFetchedWhenDynamicPathRequested(t *testing.T) {
	type pathCase struct {
		RequestedRoute   Route
		WantedParamName  string
		WantedParamValue string
	}

	casesList := []pathCase{
		{
			RequestedRoute:   NewRoute("/news/test-number-one-2020-01", http.MethodGet, handlerMock),
			WantedParamName:  "title",
			WantedParamValue: "test-number-one-2020-01",
		},
		{
			RequestedRoute:   NewRoute("/news/test-number-two-2020-05/comments", http.MethodGet, handlerMock),
			WantedParamName:  "title",
			WantedParamValue: "test-number-two-2020-05",
		},
		{
			RequestedRoute:   NewRoute("/news/test-number-three-2020-11/comments/3", http.MethodGet, handlerMock),
			WantedParamName:  "title",
			WantedParamValue: "test-number-three-2020-11",
		},
		{
			RequestedRoute:   NewRoute("/news/test-number-two-2020-05/comments/3", http.MethodGet, handlerMock),
			WantedParamName:  "id",
			WantedParamValue: "3",
		},
	}

	createRoute := func(path string, method string) Route {
		return NewRoute(path, method, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			encodedData, _ := json.Marshal(GetParamsFromContext(request.Context()))
			_, _ = writer.Write(encodedData)
		})
	}

	routesList := []Route{
		createRoute("/", http.MethodGet),
		createRoute("/", http.MethodPost),
		createRoute("/new", http.MethodGet),
		createRoute("/news", http.MethodGet),
		createRoute("/news", http.MethodPost),
		createRoute("/news/:title", http.MethodGet),
		createRoute("/news/:title/comments", http.MethodGet),
		createRoute("/news/:title/comments/:id", http.MethodGet),
	}

	router := NewRouter()
	router.AddRoutes(routesList)

	for _, data := range casesList {
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(data.RequestedRoute.GetMethod(), data.RequestedRoute.GetPath(), nil)
		router.ServeHTTP(recorder, request)

		var params PathParams
		err := json.NewDecoder(recorder.Body).Decode(&params)

		if nil != err {
			t.Errorf("router did not return params")
		} else if params.GetByName(data.WantedParamName) != data.WantedParamValue {
			t.Errorf("router returned invalid param value")
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

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(wantedRoute.GetMethod(), wantedRoute.GetPath(), nil)
	router.ServeHTTP(recorder, request)
	if http.StatusNotFound != recorder.Code {
		t.Errorf("got non-existent route")
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

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(wantedRoute.GetMethod(), wantedRoute.GetPath(), nil)
	router.ServeHTTP(recorder, request)
	if http.StatusMethodNotAllowed != recorder.Code {
		t.Errorf("got route with non-existent request method")
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
