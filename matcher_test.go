package mux

import (
	"net/http"
	"net/url"
	"testing"
)

func TestSchemeMatcher(t *testing.T) {
	schemes := []string{"http", "https"}
	matcher := newSchemeMatcher("https", "HTTP")
	request := &http.Request{
		URL: &url.URL{},
	}

	for _, v := range schemes {
		request.URL.Scheme = v
		if !matcher.Match(request) {
			t.Errorf("Scheme not matched (%v)", v)
		}
	}
}

func TestSchemeMatcherFail(t *testing.T) {
	schemes := []string{"http", "https"}
	matcher := newSchemeMatcher("httpss")
	request := &http.Request{
		URL: &url.URL{},
	}

	for _, v := range schemes {
		request.URL.Scheme = v
		if matcher.Match(request) {
			t.Errorf("Scheme matched (%v)", v)
		}
	}
}

func TestPathMatcher(t *testing.T) {
	matcher := pathMatcher("/api/v2")
	request := &http.Request{
		URL: &url.URL{
			Path: "/api/v2",
		},
	}

	if !matcher.Match(request) {
		t.Errorf("Unexpected not matched path ")
	}
}

func TestPathMatcherFail(t *testing.T) {
	matcher := pathMatcher("/api/v2")
	request := &http.Request{
		URL: &url.URL{
			Path: "/api/v1",
		},
	}

	if matcher.Match(request) {
		t.Errorf("Unexpected matched path")
	}
}

func TestMatcherFunc(t *testing.T) {
	matcherFunc := func(*http.Request) bool {
		return true
	}

	matcher := MatcherFunc(matcherFunc)

	if !matcher.Match(&http.Request{}) {
		t.Errorf("Unexpected not matched")
	}
}

func TestMatcherFuncFail(t *testing.T) {
	matcherFunc := func(*http.Request) bool {
		return false
	}

	matcher := MatcherFunc(matcherFunc)

	if matcher.Match(&http.Request{}) {
		t.Errorf("Unexpected matched")
	}
}

func TestHeaderMatcher(t *testing.T) {

	tests := []struct {
		m     func(pairs ...string) (Matcher, error)
		req   func() *http.Request
		pairs []string
	}{
		{
			m: func(pairs ...string) (Matcher, error) {
				return newHeaderMatcher(pairs...)
			},
			req: func() *http.Request {

				req := &http.Request{
					Header: http.Header{},
				}
				req.Header.Add("content-type", "applcation/json")

				return req
			},
			pairs: []string{"content-type", "applcation/json"},
		},
		{
			m: func(pairs ...string) (Matcher, error) {
				return newHeaderRegexMatcher(pairs...)
			},
			req: func() *http.Request {

				req := &http.Request{
					Header: http.Header{},
				}
				req.Header.Add("content-type", "applcation/json")

				return req
			},
			pairs: []string{"content-type", "applcation/(json|html)"},
		},
	}

	for _, test := range tests {
		matcher, err := test.m()

		if err != nil {
			t.Errorf("Unexpected error (%s)", err.Error())
		}

		req := test.req()
		if !matcher.Match(req) {
			t.Errorf("Unexpected not matched (%v)", req.Header)
		}
	}
}
