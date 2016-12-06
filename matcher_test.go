package mux

import (
	"fmt"
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

func BenchmarkSchemeMatcher(b *testing.B) {
	matcher := newSchemeMatcher("https", "http", "HTTP", "HTTPS")
	request := &http.Request{
		URL: &url.URL{
			Scheme: "HTTPS",
		},
	}

	for n := 0; n < b.N; n++ {
		matcher.Match(request)
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

func BenchmarkPathMatcher(b *testing.B) {
	matcher := pathMatcher("/api/user/2/article/4/comment/8")
	request := &http.Request{
		URL: &url.URL{
			Path: "/api/user/2/article/4/comment/8",
		},
	}

	for n := 0; n < b.N; n++ {
		matcher.Match(request)
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

func TestPathVarsMatcher(t *testing.T) {

	paths := []struct {
		pathRaw     string
		pathToMatch string
	}{
		{
			pathToMatch: "/user/:number",
			pathRaw:     "/user/1",
		},
		{
			pathToMatch: "/user/:number/comment/:number",
			pathRaw:     "/user/1/comment/99",
		},
		{
			pathToMatch: "/article/:string",
			pathRaw:     "/article/golang",
		},
		{
			pathToMatch: "/article/:string/comment/:number/subcomment/:number",
			pathRaw:     "/article/golang/comment/4/subcomment/5",
		},
		{
			pathToMatch: "/:number/:number/:number/:number/:number/:number/:number/:number/:number/:number",
			pathRaw:     "/1/1/1/1/1/1/1/1/1/1",
		},
		{
			pathToMatch: "/:string/:number/:string/:number/:string/:number/:string/:number/:string/:number",
			pathRaw:     "/dummy/1/dummy/1/dummy/1/dummy/1/dummy/1",
		},
	}

	for _, path := range paths {
		t.Run(fmt.Sprintf("Path raw: %s, Path to match %s", path.pathRaw, path.pathToMatch), func(t *testing.T) {
			matcher := newPathWithVarsMatcher(path.pathToMatch)
			request := &http.Request{
				URL: &url.URL{
					Path: path.pathRaw,
				},
			}

			if !matcher.Match(request) {
				t.Errorf("Unexpected not matched path ")
			}
		})
	}
}

func BenchmarkPathVarsMatcher(b *testing.B) {
	matcher := pathMatcher("/:string/:number/:string/:number/:string/:number/:string/:number/:string/:number")
	request := &http.Request{
		URL: &url.URL{
			Path: "/dummy/1/dummy/1/dummy/1/dummy/1/dummy/1",
		},
	}

	for n := 0; n < b.N; n++ {
		matcher.Match(request)
	}
}

func TestPathVarsMatcherFail(t *testing.T) {
	matcher := pathMatcher("/api/:number")
	request := &http.Request{
		URL: &url.URL{
			Path: "/api/echo",
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
		title   string
		m       func(pairs ...string) (Matcher, error)
		request func() *http.Request
		pairs   []string
	}{
		{
			title: "Test header match",
			m: func(pairs ...string) (Matcher, error) {
				return newHeaderMatcher(pairs...)
			},
			request: func() *http.Request {

				request := &http.Request{
					Header: http.Header{},
				}
				request.Header.Add("content-type", "applcation/json")

				return request
			},
			pairs: []string{"content-type", "applcation/json"},
		},
		{
			title: "Test regex header match",
			m: func(pairs ...string) (Matcher, error) {
				return newHeaderRegexMatcher(pairs...)
			},
			request: func() *http.Request {

				request := &http.Request{
					Header: http.Header{},
				}
				request.Header.Add("content-type", "applcation/json")

				return request
			},
			pairs: []string{"content-type", "applcation/(json|html)"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test kind: %s", test.title), func(t *testing.T) {
			matcher, err := test.m()

			if err != nil {
				t.Errorf("Unexpected error (%s)", err.Error())
			}

			request := test.request()
			if !matcher.Match(request) {
				t.Errorf("Unexpected not matched (%v)", request.Header)
			}
		})
	}
}

func BenchmarkHeaderMatcher(b *testing.B) {
	matcher, _ := newHeaderMatcher("content-type", "applcation/json")
	request := &http.Request{
		Header: http.Header{},
	}

	headers := map[string][]string{
		"Accept": {
			"text/plain",
			"text/html",
		},
		"content-type": {
			"applcation/json",
		},
		"Accept-Charset": {
			"utf-8",
		},
		"Accept-Encoding": {
			"gzip",
			"deflate",
		},
		"Accept-Language": {
			"en-US",
		},
		"Cache-Control": {
			"no-cache",
		},
		"Date": {
			"Date: Tue, 15 Nov 1994 08:12:31 GMT",
		},
		"Max-Forwards": {
			"10",
		},
	}

	for k, v := range headers {
		for _, vv := range v {
			request.Header.Add(k, vv)
		}
	}

	for n := 0; n < b.N; n++ {
		matcher.Match(request)
	}
}
