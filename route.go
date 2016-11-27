package mux

import (
	"fmt"
	"net/http"
)

// Route stores information to match a request and build URLs.
type Route struct {
	// Request handler for the route.
	handler http.Handler
	// List of matchers.
	matchers []Matcher
	// The name used to build URLs.
	name string
	// Error resulted from building a route.
	err error
	// MethodName used to build proper error messages
	methodName string
	// path used to build proper error messages
	path string

	router *Router

	buildVarsFunc BuildVarsFunc
}

//BadRouteError creates error for a bad route
type BadRouteError struct {
	r *Route
	s string
}

func newBadRouteError(r *Route, s string) *BadRouteError {
	return &BadRouteError{
		r: r,
		s: s,
	}
}

func (bre BadRouteError) Error() string {
	return fmt.Sprintf("Route -> Method: %s Path: %s Error: %s", bre.r.methodName, bre.r.path, bre.s)
}

// Match matches the route against the request.
func (r *Route) triggerMatching(req *http.Request) *Route {
	if r.err != nil {
		return nil
	}

	// Match everything.
	for _, m := range r.matchers {
		if matched := m.Match(req); !matched {
			return nil
		}
	}

	return r
}

// GetError returns an error resulted from building the route, if any.
func (r *Route) GetError() error {
	return r.err
}

// HasError check if an error exists.
func (r *Route) HasError() bool {
	if r.err == nil {
		return false
	}

	return true
}

// Handler sets a handler for the route.
func (r *Route) Handler(handler http.Handler) *Route {
	if r.err == nil {
		r.handler = handler
	}
	return r
}

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(handler))
}

// GetHandler returns the handler for the route, if any.
func (r *Route) GetHandler() http.Handler {
	return r.handler
}

// Name sets the name for the route, used to build URLs.
func (r *Route) Name(name string) *Route {

	if r.name != "" {
		r.err = newBadRouteError(r, fmt.Sprintf("route already has name %q, can't set %q", r.name, name))
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
func (r *Route) addMatcher(m Matcher) *Route {
	if r.err == nil {
		r.matchers = append(r.matchers, m)
	}
	return r
}

// Path adds a matcher for the URL path.
// It accepts a path with zero variables. The
// template must start with a "/".
// For example:
//
//     r := mux.NewRouter()
//     r.Path("/billing/").Handler(BillingHandler)
//
func (r *Route) Path(path string) *Route {

	if r.path != "" {
		r.err = newBadRouteError(r, fmt.Sprintf("route already has path can't set a new path %v", path))
	}

	r.path = path
	r.addMatcher(pathMatcher(path))

	return r
}

// Schemes adds a matcher for URL schemes.
// It accepts a sequence of schemes to be matched, e.g.: "http", "https".
func (r *Route) Schemes(schemes ...string) *Route {
	return r.addMatcher(newSchemeMatcher(schemes...))
}

// Headers adds a matcher for request header values.
// It accepts a sequence of key/value pairs to be matched. For example:
//
//     r := mux.NewRouter()
//     r.Headers("Content-Type", "application/json",
//               "X-Requested-With", "XMLHttpRequest")
//
// The above route will only match if both request header values match.
// If the value is an empty string, it will match any value if the key is set.
func (r *Route) Headers(pairs ...string) *Route {
	if r.err == nil {
		var headers map[string]string
		headers, r.err = convertStringsToMap(isEvenPairs, pairs...)
		return r.addMatcher(headerMatcher(headers))
	}
	return r
}

// MatcherFunc adds a custom function to be used as request matcher.
func (r *Route) MatcherFunc(f MatcherFunc) *Route {
	return r.addMatcher(f)
}

// BuildVarsFunc is the function signature used by custom build variable
// functions (which can modify route variables before a route's URL is built).
type BuildVarsFunc func(map[string]string) map[string]string

// BuildVarsFunc adds a custom function to be used to modify build variables
// before a route's URL is built.
func (r *Route) BuildVarsFunc(f BuildVarsFunc) *Route {
	r.buildVarsFunc = f
	return r
}

// prepareVars converts the route variable pairs into a map. If the route has a
// BuildVarsFunc, it is invoked.
func (r *Route) prepareVars(pairs ...string) (map[string]string, error) {
	m, err := convertStringsToMap(isEvenPairs, pairs...)
	if err != nil {
		return nil, err
	}
	return r.buildVars(m), nil
}

func (r *Route) buildVars(m map[string]string) map[string]string {
	if r.buildVarsFunc != nil {
		m = r.buildVarsFunc(m)
	}
	return m
}
