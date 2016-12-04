package mux

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	kindNormalPath = iota
	kindVarsPath
	kindRegexPath
)

// RouteInterface that you can create your own custom route
type RouteInterface interface {
	HasVars() bool
	HasError() bool
	SetError(error)
	GetError() error
	HasHandler() bool
	GetHandler() http.Handler
	Handler(http.Handler)
	SetMethodName(string)
	GetMethodName() string
	ExtractVars(req *http.Request) Vars
	GetPath() string
	Path(string) RouteInterface
	HandlerFunc(handler func(http.ResponseWriter, *http.Request)) RouteInterface
	GetMatchers() Matchers
	Kind() int
	Match(req *http.Request) RouteInterface
}

// Route stores information to match a request and build URLs.
type Route struct {
	// kind of route (regex, vars or normal)
	kind int
	// Request handler for the route.
	handler http.Handler
	// List of Matchers.
	ms Matchers
	// The name used to build URLs.
	name string
	// Error resulted from building a route.
	err error
	// MethodName used to build proper error messages
	methodName string
	// path used to build proper error messages
	path string
	// varIndexies used to extract vars
	varIndexies map[string]int

	router *Router
}

// NewRouter returns a new route instance.
func newRoute(router *Router) RouteInterface {
	return &Route{
		router:      router,
		ms:          Matchers([]Matcher{}),
		varIndexies: map[string]int{},
	}
}

// Match matches the route against the request.
func (r *Route) Match(req *http.Request) RouteInterface {
	if r.err != nil {
		return nil
	}

	// Match everything.
	for _, m := range r.ms {
		if matched := m.Match(req); !matched {
			return nil
		}
	}

	return r
}

// HasHandler returns ture if route has a handler.
func (r *Route) HasHandler() bool {
	if r.handler == nil {
		return false
	}
	return true
}

// GetHandler returns the handler for the route, if any.
func (r *Route) GetHandler() http.Handler {
	return r.handler
}

// SetHandler sets a handler for the route.
func (r *Route) Handler(h http.Handler) {
	if r.err == nil {
		r.handler = h
	}
}

// HasError check if an error exists.
func (r *Route) HasError() bool {
	if r.err == nil {
		return false
	}

	return true
}

// GetError returns an error resulted from building the route, if any.
func (r *Route) GetError() error {
	return r.err
}

// SetError set an error if any.
func (r *Route) SetError(err error) {
	r.err = err
}

//SetMethodName set the method name for the route
func (r *Route) SetMethodName(m string) {
	r.methodName = m
}

// GetMethodName get the method name for the route
func (r *Route) GetMethodName() string {
	return r.methodName
}

// GetMatchers get the Matchers for the route
func (r *Route) GetMatchers() Matchers {
	return r.ms
}

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	r.Handler(http.HandlerFunc(handler))
	return r
}

// Name sets the name for the route, used to build URLs.
func (r *Route) Name(name string) *Route {

	if r.name != "" {
		r.err = NewBadRouteError(r, fmt.Sprintf("route already has name %q, can't set %q", r.name, name))
		return r
	}

	if r.err == nil {
		r.name = name
	}

	return r
}

// GetName returns the name for the route, if any.
func (r *Route) GetName() string {
	return r.name
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m Matcher) RouteInterface {
	if r.err == nil {
		r.ms = append(r.ms, m)
	}
	return r
}

// Path adds a matcher for the URL path.
// It accepts a path with zero variables. The
// template must start with a "/".
// For example:
//
//     r := mux.Classic()
//     r.Path("/billing/").Handler(BillingHandler)
//
func (r *Route) Path(path string) RouteInterface {

	if r.path != "" {
		r.err = NewBadRouteError(r, fmt.Sprintf("route already has path can't set a new path %v", path))
	}

	var matcher Matcher
	switch {
	case containsRegex(path):
		matcher = newPathRegexMatcher(path)
		r.kind = kindRegexPath
	case containsVars(path):
		matcher = newPathWithVarsMatcher(path)
		r.extractVarsIndexies(":", path)
		r.kind = kindVarsPath
	default:
		matcher = pathMatcher(path)
		r.kind = kindNormalPath
	}

	r.path = path
	r.addMatcher(matcher)

	return r
}

//GetPath returns the handler for the route, if any.
func (r *Route) GetPath() string {
	return r.path
}

func (r *Route) extractVarsIndexies(prefix string, path string) {

	urlSeg := strings.Split(path, "/")

	indexies := map[string]int{}
	var count int
	for k, v := range urlSeg {
		if strings.HasPrefix(v, prefix) {

			if _, found := indexies[v]; !found {
				indexies[v] = k
				continue
			}

			count++
			indexies[v+string(count)] = k
		}
	}

	r.varIndexies = indexies
}

//HasVars check if path has any vars
func (r *Route) HasVars() bool {
	if 0 == len(r.varIndexies) {
		return false
	}
	return true
}

type Vars map[string]string

// Get return the key value, of the current *http.Request queries
func (v Vars) Get(key string) string {
	if value, found := v[key]; found {
		return value
	}
	return ""
}

// GetAll returns all queries of the current *http.Request queries
func (v Vars) GetAll() map[string]string {
	return v
}

//ExtractVars extract all vars of the current path
func (r *Route) ExtractVars(req *http.Request) Vars {

	urlSeg := strings.Split(req.URL.Path, "/")

	vars := Vars(map[string]string{})

	for k, v := range r.varIndexies {
		vars[k] = urlSeg[v]
	}

	return vars
}

// Schemes adds a matcher for URL schemes.
// It accepts a sequence of schemes to be matched, e.g.: "http", "https".
func (r *Route) Schemes(schemes ...string) RouteInterface {
	return r.addMatcher(newSchemeMatcher(schemes...))
}

// Headers adds a matcher for request header values.
// It accepts a sequence of key/value pairs to be matched. For example:
//
//     r := mux.Classic()
//     r.Headers("Content-Type", "application/json",
//               "X-Requested-With", "XMLHttpRequest")
//
// The above route will only match if both request header values match.
// If the value is an empty string, it will match any value if the key is set.
func (r *Route) Headers(pairs ...string) RouteInterface {
	if r.err != nil {
		return r
	}

	matcher, err := newHeaderMatcher(pairs...)

	if err != nil {
		r.err = err
	}

	r.addMatcher(matcher)

	return r
}

// HeadersRegex adds a matcher for request header values.
// It accepts a sequence of key/value pairs to be matched. For example:
//
//     r := mux.Classic()
//     r.Headers("Content-Type", "application/(json|html)")
//
// The above route will only match if both request header values match.
// If the value is an empty string, it will match any value if the key is set.
func (r *Route) HeadersRegex(pairs ...string) RouteInterface {
	if r.err != nil {
		return r
	}

	matcher, err := newHeaderRegexMatcher(pairs...)

	if err != nil {
		r.err = err
	}

	r.addMatcher(matcher)

	return r
}

// MatcherFunc adds a custom function to be used as request matcher.
func (r *Route) MatcherFunc(f MatcherFunc) RouteInterface {
	return r.addMatcher(f)
}

//Kind returns kind of route
func (r *Route) Kind() int {
	return r.kind
}
