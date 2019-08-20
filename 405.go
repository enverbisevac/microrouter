package microrouter

import (
	"io"
	"net/http"
)

func Http405Text(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(405)
	io.WriteString(w, "Method not allowed.")
}

func Http405Html(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(405)
	io.WriteString(w, "<h1>Method not allowed.</h1>")
}
