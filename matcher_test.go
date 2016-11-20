package mux

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMethodMatcher(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH"}
	matcher := newMethodMatcher("GET", "post", "Put", "PaTch")
	request := &http.Request{}

	for _, v := range methods {
		request.Method = v
		if !matcher.Match(request) {
			t.Errorf("Method not matched (%v)", v)
		}
	}
}

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
