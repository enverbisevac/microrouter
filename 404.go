package microrouter

import (
	"io"
	"net/http"
)

func Http404Text(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(404)
	_, _ = io.WriteString(w, "404 - Page not found.")
}

func Http404Html(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(404)
	_, _ = io.WriteString(w, "<h1>404 - Page not found.</h1>")
}
