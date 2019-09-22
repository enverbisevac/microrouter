package microrouter

import (
	"github.com/gorilla/reverse"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync/atomic"
)

const (
	defaultContentType = "text/html"
	textContentType    = "text/plain"
)

const (
	methodNotFound uint8 = iota
	pathNotFound
	methodAndPathFound
)

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

type route struct {
	pattern     string
	handlerFunc http.HandlerFunc
	counter     uint64
}

type RegexHandler interface {
	http.Handler
	SetRootGroup(group *Group)
	SetNotFoundHandler(contentType string, handlerFunc http.HandlerFunc)
	SetMethodNotFoundHandler(contentType string, handlerFunc http.HandlerFunc)
	Reverse(name string, values url.Values) (string, error)
}

type regexResolver struct {
	routes         []route
	names          map[string]string
	cache          map[string]*reverse.Regexp
	notFound       map[string]http.HandlerFunc // Content type based not found resource
	methodNotFound map[string]http.HandlerFunc
	visited        chan<- string
}

func newRegexResolver() *regexResolver {
	notFound := make(map[string]http.HandlerFunc)
	notFound[defaultContentType] = Http404Html
	notFound[textContentType] = Http404Text
	methodNotFound := make(map[string]http.HandlerFunc)
	methodNotFound[defaultContentType] = Http405Html
	methodNotFound[textContentType] = Http405Text
	visited := make(chan string)

	regex := &regexResolver{
		visited:        visited,
		names:          make(map[string]string),
		cache:          make(map[string]*reverse.Regexp),
		notFound:       notFound,
		methodNotFound: methodNotFound,
	}
	// run routine to calculate traffic
	go regex.visitedFunc(visited)
	return regex
}

func (r *regexResolver) visitedFunc(visited <-chan string) {
	for {
		val := <-visited
		for _, route := range r.routes {
			if route.pattern == val {
				log.Printf("Url %s have %d hits", val, atomic.LoadUint64(&route.counter))
				break
			}
		}
	}
}

func (r *regexResolver) SetRootGroup(group *Group) {
	// generate all urls
	group.generate("")
	// copy to routes
	for _, routeObject := range group.routes {
		log.Println(routeObject.pattern)
		r.routes = append(r.routes, route{
			pattern:     routeObject.pattern,
			handlerFunc: routeObject.handlerFunc,
		})
		r.names[routeObject.pattern] = group.names[routeObject.pattern]
		cache, err := reverse.CompileRegexp(routeObject.pattern)
		if err != nil {
			continue
		}
		// set cache on compiled regex
		r.cache[routeObject.pattern] = cache
	}
}

func (r *regexResolver) ReverseWithMethod(name, method string, values url.Values) (string, error) {
	pattern := r.names[name]
	// initial method argument in pattern can be any of this values
	revertValues := url.Values{
		"": {method},
	}
	for key, value := range values {
		for _, text := range value {
			revertValues.Add(key, text)
		}
	}
	result, err := r.cache[pattern].Revert(revertValues)
	if err != nil {
		return "", err
	}
	// result will be for example GET /article/1
	// we need only second element after splitting
	return strings.Split(result, " ")[1], nil
}

func (r *regexResolver) Reverse(name string, values url.Values) (string, error) {
	return r.ReverseWithMethod(name, "GET", values)
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
	for index, route := range r.routes {
		if r.cache[route.pattern].MatchString(check) {
			route.handlerFunc(res, req)
			atomic.AddUint64(&r.routes[index].counter, 1)
			r.visited <- route.pattern
			return
		}
	}

	// following code shoud not be focused too much on speed
	// bcoz error 405 or 404 are very rare
	result := pathNotFound
	for _, route := range r.routes {
		if checkMethod(req.Method, req.URL.Path, route.pattern) == methodNotFound {
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

func getContentType(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		contentType = defaultContentType
	}
	return contentType
}

func checkMethod(inputMethod, inputPath, pattern string) uint8 {
	log.Printf("Checking request method %s with pattern %s", inputMethod, pattern)
	if inputMethod == "OPTIONS" {
		log.Println("OPTIONS must pass check")
		return methodAndPathFound
	}
	splitter := strings.Split(pattern, " ")
	// we have to check path of our request url
	path := splitter[1]
	regex, _ := regexp.Compile(path)
	if !regex.MatchString(inputPath) {
		return pathNotFound
	}
	// check request method if path was founded
	method := splitter[0]
	regex, _ = regexp.Compile(method)
	if !regex.MatchString(inputMethod) {
		return methodNotFound
	}
	return methodAndPathFound
}
