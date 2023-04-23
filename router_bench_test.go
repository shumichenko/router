package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkRouter_ServeHTTP(b *testing.B) {
	var routesList []Route

	methodsList := []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	dictionary := "abc1def2ghi3jkl4mno5pqr6stu7vwx8yz9_0-"
	for _, method := range methodsList {
		for _, symbol := range dictionary {
			routesList = append(
				routesList,
				NewRoute("/"+string(symbol)+"test-example_high/load", method, handlerMock),
				NewRoute("/"+string(symbol)+"test-example_high/load/route/test/help", method, handlerMock),
				NewRoute("/"+string(symbol)+"test-example_high/load/route/test/help/test", method, handlerMock),
				NewRoute("/"+string(symbol)+"test-example_high/load/route/test/exam/test/a", method, handlerMock),
			)
		}
	}
	routesList = append(
		routesList,
		NewRoute("/news/:id/comments/:id/statistics/test/test", http.MethodGet, handlerMock),
	)

	router := NewRouter()
	router.AddRoutes(routesList)

	b.ReportAllocs()
	b.ResetTimer()

	recorder := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		request, _ := http.NewRequest(http.MethodGet, "/news/103/comments/10/statistics/test/test", nil)
		router.ServeHTTP(recorder, request)
	}
}
