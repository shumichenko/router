package router

import (
	"net/http"
	"testing"
)

func TestPanicWhenEmptyPathProvided(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("route was created with empty path")
		}
	}()
	NewRoute("", http.MethodGet, handlerMock)
}

func TestPanicWhenEmptyMethodProvided(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("route was created with empty method")
		}
	}()
	NewRoute("/", "", handlerMock)
}

func TestPanicWhenPathWithNoStartingSlashProvided(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("route was created without starting slash")
		}
	}()
	NewRoute("news", "", handlerMock)
}

func TestPanicWhenRestrictedCharsInPath(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("route was created with restricted chars")
		}
	}()
	NewRoute("/news?", "", handlerMock)
}

func TestTrailingSlashRemoved(t *testing.T) {
	defer func() {
		if r := recover(); nil == r {
			t.Errorf("route was created with trailing slash")
		}
	}()
	NewRoute("/news/", http.MethodGet, handlerMock)
}
