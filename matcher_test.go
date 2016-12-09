package mux

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

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

func TestPathMatchers(t *testing.T) {

	tests := []struct {
		title        string
		pathRaw      string
		pathToMatch  string
		buildMatcher func(string) Matcher
	}{
		{
			title:       "Path noraml matcher ",
			pathToMatch: "/api/echo",
			pathRaw:     "/api/echo",
			buildMatcher: func(path string) Matcher {
				return pathMatcher(path)
			},
		},
		{
			title:       "Path vars matcher (one number var segment)",
			pathToMatch: "/user/:number",
			pathRaw:     "/user/1",
			buildMatcher: func(path string) Matcher {
				return newPathWithVarsMatcher(path)
			},
		},
		{
			title:       "Path vars matcher (two number var segments)",
			pathToMatch: "/user/:number/comment/:number",
			pathRaw:     "/user/1/comment/99",
			buildMatcher: func(path string) Matcher {
				return newPathWithVarsMatcher(path)
			},
		},
		{
			title:       "Path vars matcher (two string var segments)",
			pathToMatch: "/article/:string",
			pathRaw:     "/article/golang",
			buildMatcher: func(path string) Matcher {
				return newPathWithVarsMatcher(path)
			},
		},
		{
			title:       "Path vars matcher (one number and one string var segment)",
			pathToMatch: "/article/:string/comment/:number/subcomment/:number",
			pathRaw:     "/article/golang/comment/4/subcomment/5",
			buildMatcher: func(path string) Matcher {
				return newPathWithVarsMatcher(path)
			},
		},
		{
			title:       "Path vars matcher (many number var segments)",
			pathToMatch: "/:number/:number/:number/:number/:number/:number/:number/:number/:number/:number",
			pathRaw:     "/1/1/1/1/1/1/1/1/1/1",
			buildMatcher: func(path string) Matcher {
				return newPathWithVarsMatcher(path)
			},
		},
		{
			title:       "Path vars matcher (many number and string var segments)",
			pathToMatch: "/:string/:number/:string/:number/:string/:number/:string/:number/:string/:number",
			pathRaw:     "/dummy/1/dummy/1/dummy/1/dummy/1/dummy/1",
			buildMatcher: func(path string) Matcher {
				return newPathWithVarsMatcher(path)
			},
		},
		{
			title:       "Path regex matcher (many regex segments)",
			pathToMatch: "/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}",
			pathRaw:     "/dummy/1/dummy/1/dummy/1/dummy/1/dummy/1",
			buildMatcher: func(path string) Matcher {
				return newPathRegexMatcher(path)
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Path raw: %s, Path to match %s "+test.title, test.pathRaw, test.pathToMatch), func(t *testing.T) {
			matcher := test.buildMatcher(test.pathToMatch)
			request := &http.Request{
				URL: &url.URL{
					Path: test.pathRaw,
				},
			}

			if !matcher.Match(request) {
				t.Errorf("Unexpected not matched path ")
			}
		})
	}
}

func BenchmarkPathMatchers(b *testing.B) {

	benchmarks := []struct {
		title        string
		pathRaw      string
		pathToMatch  string
		buildMatcher func(string) Matcher
	}{
		{
			title:       "Path noraml matcher (2 URL segments)",
			pathToMatch: "/api/echo",
			pathRaw:     "/api/echo",
			buildMatcher: func(path string) Matcher {
				return pathMatcher(path)
			},
		},
		{
			title:       "Path noraml matcher (7 URL segments)",
			pathToMatch: "/api/user/2/article/4/comment/8",
			pathRaw:     "/api/user/2/article/4/comment/8",
			buildMatcher: func(path string) Matcher {
				return pathMatcher(path)
			},
		},
		{
			title:       "Path vars matcher (many vars segments)",
			pathToMatch: "/:string/:number/:string/:number/:string/:number/:string/:number/:string/:number",
			pathRaw:     "/user/1",
			buildMatcher: func(path string) Matcher {
				return newPathWithVarsMatcher(path)
			},
		},
		{
			title:       "Path regex matcher (many regex segments)",
			pathToMatch: "/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}/#([a-z]){1,}/#([0-9]){1,}",
			pathRaw:     "/dummy/1/dummy/1/dummy/1/dummy/1/dummy/1",
			buildMatcher: func(path string) Matcher {
				return newPathRegexMatcher(path)
			},
		},
	}

	for _, benchmark := range benchmarks {
		matcher := benchmark.buildMatcher(benchmark.pathToMatch)
		request := &http.Request{
			URL: &url.URL{
				Path: benchmark.pathRaw,
			},
		}
		b.Run(benchmark.title, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				matcher.Match(request)
			}
		})
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

func TestHeaderMatcher(t *testing.T) {

	tests := []struct {
		title        string
		buildMatcher func(pairs ...string) (Matcher, error)
		buildRequest func() *http.Request
		pairs        []string
	}{
		{
			title: "Test header match",
			buildMatcher: func(pairs ...string) (Matcher, error) {
				return newHeaderMatcher(pairs...)
			},
			buildRequest: func() *http.Request {

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
			buildMatcher: func(pairs ...string) (Matcher, error) {
				return newHeaderRegexMatcher(pairs...)
			},
			buildRequest: func() *http.Request {

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
			matcher, err := test.buildMatcher()

			if err != nil {
				t.Errorf("Unexpected error (%s)", err.Error())
			}

			request := test.buildRequest()
			if !matcher.Match(request) {
				t.Errorf("Unexpected not matched (%v)", request.Header)
			}
		})
	}
}

func BenchmarkHeaderMatchers(b *testing.B) {

	buildRequest := func() *http.Request {
		return &http.Request{
			Header: http.Header{},
		}
	}

	tests := []struct {
		title        string
		buildMatcher func() Matcher
	}{
		{
			title: "Benchmark: Header matcher (single value)",
			buildMatcher: func() Matcher {
				matcher, _ := newHeaderMatcher("content-type", "applcation/json")
				return matcher
			},
		},
		{
			title: "Benchmark: Header matcher (double value)",
			buildMatcher: func() Matcher {
				matcher, _ := newHeaderMatcher("content-type", "applcation/json")
				return matcher
			},
		},
		{
			title: "Benchmark: Header regex matcher (single value)",
			buildMatcher: func() Matcher {
				matcher, _ := newHeaderRegexMatcher("content-type", "applcation/(json|html)", "accept", "text/(plain|html)")
				return matcher
			},
		},
		{
			title: "Benchmark: Header regex matcher (double value)",
			buildMatcher: func() Matcher {
				matcher, _ := newHeaderRegexMatcher("content-type", "applcation/(json|html)", "accept", "text/(plain|html)")
				return matcher
			},
		},
	}

	for _, test := range tests {
		request := buildRequest()
		populateHeaderWithTestData(request)
		matcher := test.buildMatcher()
		b.Run(test.title, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				matcher.Match(request)
			}
		})
	}
}

func populateHeaderWithTestData(request *http.Request) {
	headers := map[string][]string{
		"content-type": {
			"applcation/json",
		},
		"accept-charset": {
			"utf-8",
		},
		"accept-encoding": {
			"gzip",
			"deflate",
		},
		"accept-language": {
			"en-US",
		},
		"cache-control": {
			"no-cache",
		},
		"date": {
			"Date: Tue, 15 Nov 1994 08:12:31 GMT",
		},
		"max-Forwards": {
			"10",
		},
		"accept": {
			"text/plain",
			"text/html",
		},
	}

	for k, v := range headers {
		for _, vv := range v {
			request.Header.Add(k, vv)
		}
	}
}

func BenchmarkNewPathWithVarsMatcher(b *testing.B) {

	tests := []struct {
		title string
		path  string
	}{
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (one var number value)",
			path:  "/:number",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (two var number value)",
			path:  "/:number/:numnber",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (three var number value)",
			path:  "/:number/:number/:number",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (four var number value)",
			path:  "/:number/:number/:number/:number",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (five var number value)",
			path:  "/:number/:number/:number/:number/:number",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (one var string value)",
			path:  "/:string",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (two var string value)",
			path:  "/:string/:string",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (three var string value)",
			path:  "/:string/:string/:string",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (four var string value)",
			path:  "/:string/:string/:string/:string",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (five var string value)",
			path:  "/:string/:string/:string/:string/:string",
		},
		{
			title: "Benchmark: constructor of pathWithVarsMatcher (mixed var value)",
			path:  "/:number/:string/:number/:string/:number/:number/:string/:number/:string/:number",
		},
	}

	for _, test := range tests {
		b.Run(test.title, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				newPathWithVarsMatcher(test.path)
			}
		})
	}
}
