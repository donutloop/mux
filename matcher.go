package mux

import (
	"net/http"
	"strings"
)

// Matcher types try to match a request.
type Matcher interface {
	Match(*http.Request) bool
}

// headerMatcher matches the request against header values.
type headerMatcher map[string]string

func (m headerMatcher) Match(r *http.Request) bool {
	return matchMapWithString(m, r.Header, true)
}

// MatcherFunc is the function signature used by custom matchers.
type MatcherFunc func(*http.Request) bool

// Match returns the match for a given request.
func (m MatcherFunc) Match(r *http.Request) bool {
	return m(r)
}

// schemeMatcher matches the request against URL schemes.
type schemeMatcher map[string]struct{}

func newSchemeMatcher(schemes ...string) schemeMatcher {
	schemeMatcher := schemeMatcher{}

	for _, v := range schemes {
		schemeMatcher[strings.ToLower(v)] = struct{}{}
	}

	return schemeMatcher
}

func (m schemeMatcher) Match(r *http.Request) bool {
	if _, found := m[r.URL.Scheme]; found {
		return true
	}

	return false
}

//pathMatcher matches the request against a URL path.
type pathMatcher string

func (m pathMatcher) Match(r *http.Request) bool {
	if strings.Compare(string(m), r.URL.Path) == 0 {
		return true
	}

	return false
}

//methodMatcher matches the string against a method.
type methodMatcher map[string]struct{}

// newMatcher
func newMethodMatcher() methodMatcher {
	return methodMatcher(methods)
}

// methods all possible standard methods
var methods = map[string]struct{}{
	http.MethodGet:     struct{}{},
	http.MethodPost:    struct{}{},
	http.MethodPatch:   struct{}{},
	http.MethodDelete:  struct{}{},
	http.MethodHead:    struct{}{},
	http.MethodPut:     struct{}{},
	http.MethodOptions: struct{}{},
}

func (m methodMatcher) Match(method string) bool {

	if _, found := m[method]; found {
		return true
	}

	return false
}
