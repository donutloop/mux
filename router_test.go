package mux

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type routeTest struct {
	title      string
	path       string
	method     string
	statusCode int
	kind       string
	route      func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request))
}

func TestPath(t *testing.T) {

	tests := []routeTest{
		{
			title:      "(GET) Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Get(path, handler)
			},
		},
		{
			title:      "(POST) Path route with single path",
			path:       "/api/",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Post(path, handler)
			},
		},
		{
			title:      "(DELETE) Path route with single path",
			path:       "/api",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Delete(path, handler)
			},
		},
		{
			title:      "(PUT) Path route with single path",
			path:       "/api",
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Put(path, handler)
			},
		},
		{
			title:      "(Head) Path route with single path",
			path:       "/api/",
			method:     http.MethodHead,
			statusCode: http.StatusOK,
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Head(path, handler)
			},
		},
		{
			title:      "(Options) Path route with single path",
			path:       "/api/",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Options(path, handler)
			},
		},
		{
			title:      "(GET) Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, path, handler)
			},
		},
		{
			title:      "(POST) Path route with single path",
			path:       "/api/",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, path, handler)
			},
		},
		{
			title:      "(DELETE) Path route with single path",
			path:       "/api/",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, path, handler)
			},
		},
		{
			title:      "(PUT) Path route with single path",
			path:       "/api/",
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, path, handler)
			},
		},
		{
			title:      "(Head) Path route with single path",
			path:       "/api/",
			method:     http.MethodHead,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, path, handler)
			},
		},
		{
			title:      "(Options) Path route with single path",
			path:       "/api/",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, path, handler)
			},
		},
		{
			title:      "(GET) Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "Handler",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Handle(method, path, http.HandlerFunc(handler))
			},
		},
		{
			title:      "(POST) Path route with single path",
			path:       "/api/",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			kind:       "Handler",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Handle(method, path, http.HandlerFunc(handler))
			},
		},
		{
			title:      "(DELETE) Path route with single path",
			path:       "/api/",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
			kind:       "Handler",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Handle(method, path, http.HandlerFunc(handler))
			},
		},
		{
			title:      "(PUT) Path route with single path",
			path:       "/api/",
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			kind:       "Handler",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Handle(method, path, http.HandlerFunc(handler))
			},
		},
		{
			title:      "(Head) Path route with single path",
			path:       "/api/",
			method:     http.MethodHead,
			statusCode: http.StatusOK,
			kind:       "Handler",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Handle(method, path, http.HandlerFunc(handler))
			},
		},
		{
			title:      "(Options) Path route with single path",
			path:       "/api/",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "Handler",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Handle(method, path, http.HandlerFunc(handler))
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test: %s path: %s method %s kind: %s", test.title, test.path, test.method, test.kind), func(t *testing.T) {
			code, ok := testRoute(test)

			if !ok {
				t.Errorf("Expected status code %v, Actucal status code %v", test.statusCode, code)
			}
		})
	}
}

func testRoute(rt routeTest) (int, bool) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("succesfully"))
	}
	r := NewRouter()
	rt.route(r, rt.path, rt.method, handler)

	req, _ := http.NewRequest(rt.method, "http://localhost"+rt.path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var content bytes.Buffer
	_, err := io.Copy(&content, res.Body)

	if err != nil {
		return -1, false
	}

	if res.Code != rt.statusCode && content.String() != "succesfully" {
		return res.Code, false
	}

	return res.Code, true
}
