package router

import (
	"net/http"
	"testing"
)

func TestPanicWhenEmptyPathProvided(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("Route was created with empty path")
		}
	}()
	NewRoute("", http.MethodGet, handlerMock)
}

func TestPanicWhenEmptyMethodProvided(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("Route was created with empty method")
		}
	}()
	NewRoute("/", "", handlerMock)
}

func TestPanicWhenPathWithNoStartingSlashProvided(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("Route was created without starting slash")
		}
	}()
	NewRoute("news", "", handlerMock)
}

func TestTrailingSlashRemoved(t *testing.T) {
	route := NewRoute("/news/", http.MethodGet, handlerMock)
	if '/' == route.Path[len(route.Path)-1] {
		t.Errorf("Route was created with trailing slash")
	}
}
