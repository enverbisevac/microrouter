package microrouter

import (
	"net/http"
	"testing"
)

func TestCheckMethod(t *testing.T) {
	got := checkMethod("GET", "/hello", "GET /hello")
	exp := methodAndPathFound
	if got != exp {
		t.Errorf("Expected result %d, but we got %d", exp, got)
	}

	got = checkMethod("GET", "/hello", "POST /hello")
	exp = methodNotFound
	if got != exp {
		t.Errorf("Expected result %d, but we got %d", exp, got)
	}

	got = checkMethod("GET", "/", "POST /hello")
	exp = pathNotFound
	if got != exp {
		t.Errorf("Expected result %d, but we got %d", exp, got)
	}

	got = checkMethod("GET", "/hello", "(GET|POST) /hello")
	exp = methodAndPathFound
	if got != exp {
		t.Errorf("Expected result %d, but we got %d", exp, got)
	}
}

func TestGetContentType(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "")
	got := getContentType(req)
	exp := defaultContentType
	if got != exp {
		t.Errorf("Expected content type %s, but we got %s", exp, got)
	}

	req.Header.Set("Content-Type", "application/json")
	got = getContentType(req)
	exp = "application/json"
	if got != exp {
		t.Errorf("Expected content type %s, but we got %s", exp, got)
	}
}
