package microrouter

import (
	"net/http"
	"testing"
)

func TestHttp404Html(t *testing.T) {
	defer testHttp(t, Http404Html, http.StatusNotFound, "<h1>404 - Page not found.</h1>")
}

func TestHttp404Text(t *testing.T) {
	defer testHttp(t, Http404Text, http.StatusNotFound, "404 - Page not found.")
}
