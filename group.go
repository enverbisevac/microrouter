package microrouter

import (
	"fmt"
	"net/http"
	"strings"
)

type Grouper interface {
	Add(pattern string, handlerFunc http.HandlerFunc, methods ...string) error
	AddWithName(name, pattern string, handlerFunc http.HandlerFunc, methods ...string) error
	AddGroup(prefix string, group *Group) error
}

type Group struct {
	parent      *Group
	routes      []route
	children    map[string]*Group
	names       map[string]string
	middlewares MiddlewareChain
}

func NewGroup() *Group {
	return &Group{
		children:    make(map[string]*Group),
		names:       make(map[string]string),
		middlewares: MiddlewareChain{},
	}
}

func (group *Group) AddGroup(prefix string, newGroup *Group) error {
	newGroup.parent = group
	group.children[prefix] = newGroup
	return nil
}

func (group *Group) generate(prefix string) {
	for _prefix, child := range group.children {
		child.generate(_prefix)
	}
	for _, route := range group.routes {
		name := group.names[route.pattern]
		parts := strings.Split(route.pattern, " ")
		pattern := fmt.Sprintf("%s%s", prefix, parts[1])
		// split (GET|POST) for example
		methods := strings.Split(parts[0][1:len(parts[0])-1], "|")
		if group.parent != nil {
			err := group.parent.AddWithName(name, pattern, route.handlerFunc, methods...)
			if err != nil {
				continue
			}
		}
	}
}

func (group *Group) Add(pattern string, handlerFunc http.HandlerFunc, methods ...string) error {
	fullPattern := generateFullPattern(pattern, methods...)
	// set handler for this pattern
	group.routes = append(group.routes, route{
		pattern:     fullPattern,
		handlerFunc: handlerFunc,
	})
	return nil
}

func (group *Group) AddWithName(name, pattern string, handlerFunc http.HandlerFunc, methods ...string) error {
	err := group.Add(pattern, handlerFunc, methods...)
	group.names[name] = generateFullPattern(pattern, methods...)
	return err
}

func generateFullPattern(pattern string, methods ...string) string {
	if pattern == "" {
		pattern = "/$"
	}
	methodsString := "(GET)"
	if len(methods) > 0 {
		methodsString = fmt.Sprintf("(%s)", strings.Join(methods, "|"))
	}
	fullPattern := strings.Join([]string{methodsString, pattern}, " ")
	return fullPattern
}
