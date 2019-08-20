package microrouter

import (
	"log"
	"net/http"
	"time"
)

type loginResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoginResponseWriter(w http.ResponseWriter) *loginResponseWriter {
	return &loginResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loginResponseWriter) WriteHeader(code int) {
	lrw.ResponseWriter.WriteHeader(code)
	lrw.statusCode = code
}

func LoggerMiddleware() MiddlewareInterceptor {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		lrw := newLoginResponseWriter(w)
		next(lrw, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	}
}
