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
