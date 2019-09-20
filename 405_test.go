package microrouter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHttp(t *testing.T, handlerFunc http.HandlerFunc, code int, body string) func() {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Error msg %v", err)
		t.FailNow()
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rec, req)

	if rec.Code != code {
		t.Errorf("Status code expected %d, actual %d", code, rec.Code)
	}
	got := rec.Body.String()
	if body != got {
		t.Errorf("Response expected %s, actual %s", body, got)
	}
	return func() {}
}

func TestHttp405Html(t *testing.T) {
	defer testHttp(t, Http405Html, http.StatusMethodNotAllowed, "<h1>Method not allowed.</h1>")
}

func TestHttp405Text(t *testing.T) {
	defer testHttp(t, Http405Text, http.StatusMethodNotAllowed, "Method not allowed.")
}
