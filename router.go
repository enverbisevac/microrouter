package microrouter

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const defaultContentType = "text/html"

const (
	methodNotFound uint8 = iota
	pathNotFound
	methodAndPathFound
)

type RegexHandler interface {
	http.Handler
	Add(pattern string, handlerFunc http.HandlerFunc, methods ...string) error
	SetNotFoundHandler(contentType string, handlerFunc http.HandlerFunc)
	SetMethodNotFoundHandler(contentType string, handlerFunc http.HandlerFunc)
}

type Router struct {
	RegexHandler
	middlewares MiddlewareChain
}

func NewRouter() *Router {
	return &Router{
		RegexHandler: newRegexResolver(),
		middlewares:  MiddlewareChain{},
	}
}

func (router *Router) Use(interceptor MiddlewareInterceptor) {
	router.middlewares = append(router.middlewares, interceptor)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := router.middlewares.Handler(router.RegexHandler.ServeHTTP)
	handler.ServeHTTP(w, r)
}

type regexResolver struct {
	handlers       map[string]http.HandlerFunc
	cache          map[string]*regexp.Regexp
	notFound       map[string]http.HandlerFunc // Content type based not found resource
	methodNotFound map[string]http.HandlerFunc
	errorHandler   map[string]http.HandlerFunc // Content type based internal server error
}

func newRegexResolver() *regexResolver {
	notFound := make(map[string]http.HandlerFunc)
	notFound[defaultContentType] = http.NotFound
	notFound["text/plain"] = http.NotFound
	methodNotFound := make(map[string]http.HandlerFunc)
	methodNotFound[defaultContentType] = Http405Html
	methodNotFound["text/plain"] = Http405Text
	return &regexResolver{
		handlers:       make(map[string]http.HandlerFunc),
		cache:          make(map[string]*regexp.Regexp),
		notFound:       notFound,
		methodNotFound: methodNotFound,
		errorHandler:   make(map[string]http.HandlerFunc),
	}
}

func (r *regexResolver) Add(pattern string, handlerFunc http.HandlerFunc, methods ...string) error {
	methodsString := "GET"
	if len(methods) > 0 {
		methodsString = fmt.Sprintf("(%s)", strings.Join(methods, "|"))
	}
	fullPattern := strings.Join([]string{methodsString, pattern}, " ")
	r.handlers[fullPattern] = handlerFunc
	cache, err := regexp.Compile(fullPattern)
	if err != nil {
		return err
	}
	r.cache[fullPattern] = cache
	return nil
}

func (r *regexResolver) SetNotFoundHandler(contentType string, handlerFunc http.HandlerFunc) {
	r.notFound[contentType] = handlerFunc
}

func (r *regexResolver) SetMethodNotFoundHandler(contentType string, handlerFunc http.HandlerFunc) {
	r.methodNotFound[contentType] = handlerFunc
}

func (r *regexResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// what to check for example GET /about
	check := strings.Join([]string{req.Method, req.URL.Path}, " ")
	// try to find in routes table
	for pattern, handlerFunc := range r.handlers {
		if r.cache[pattern].MatchString(check) == true {
			handlerFunc(res, req)
			return
		}
	}

	// following code shoud not be focused too much on speed
	// bcoz error 405 or 404 are very rare
	result := pathNotFound
	for pattern := range r.handlers {
		if checkMethod(req.Method, req.URL.Path, pattern) == methodNotFound {
			result = methodNotFound
			break
		}
	}

	if result == methodNotFound {
		r.Http405(res, req)
		return
	}

	// resource not found
	r.Http404(res, req)
}

func (r *regexResolver) Http404(res http.ResponseWriter, req *http.Request) {
	contentType := getContentType(req)
	r.notFound[contentType](res, req)
}

func (r *regexResolver) Http405(res http.ResponseWriter, req *http.Request) {
	contentType := getContentType(req)
	r.methodNotFound[contentType](res, req)
}

func (r *regexResolver) Http500(res http.ResponseWriter, req *http.Request) {
	contentType := getContentType(req)
	r.errorHandler[contentType](res, req)
}

func getContentType(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		contentType = defaultContentType
	}
	return contentType
}

func checkMethod(inputMethod, inputPath, pattern string) uint8 {
	log.Printf("Checking request method %s with pattern %s", inputMethod, pattern)
	splitter := strings.Split(pattern, " ")
	// we have to check path of our request url
	path := splitter[1]
	regex, _ := regexp.Compile(path)
	if regex.MatchString(inputPath) != true {
		return pathNotFound
	}
	// check request method if path was founded
	method := splitter[0]
	regex, _ = regexp.Compile(method)
	if regex.MatchString(inputMethod) != true {
		return methodNotFound
	}
	return methodAndPathFound
}
