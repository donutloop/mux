package mux

// matcher types try to match a request.
import (
	"net/http"
	"strings"
)

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

// methodMatcher matches the request against HTTP methods.
type methodMatcher map[string]struct{}

func newMethodMatcher(methods ...string) methodMatcher {
	methodMatcher := methodMatcher{}

	for _, v := range methods {
		methodMatcher[strings.ToUpper(v)] = struct{}{}
	}

	return methodMatcher
}

func (m methodMatcher) Match(r *http.Request) bool {
	if _, found := m[r.Method]; found {
		return true
	}

	return false
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
