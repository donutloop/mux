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
}

func TestPath(t *testing.T) {

	tests := []routeTest{
		{
			title:      "(GET) Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
		{
			title:      "(POST) Path route with single path",
			path:       "/api/",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
		},
		{
			title:      "(DELETE) Path route with single path",
			path:       "/api",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
		},
		{
			title:      "(PUT) Path route with single path",
			path:       "/api",
			method:     http.MethodPut,
			statusCode: http.StatusOK,
		},
		{
			title:      "(Head) Path route with single path",
			path:       "/api/",
			method:     http.MethodHead,
			statusCode: http.StatusOK,
		},
		{
			title:      "(Options) Path route with single path",
			path:       "/api/",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
		},
		{
			title:      "(GET) Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
		},
		{
			title:      "(POST) Path route with single path",
			path:       "/api/",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
		},
		{
			title:      "(DELETE) Path route with single path",
			path:       "/api/",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
		},
		{
			title:      "(PUT) Path route with single path",
			path:       "/api/",
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
		},
		{
			title:      "(Head) Path route with single path",
			path:       "/api/",
			method:     http.MethodHead,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
		},
		{
			title:      "(Options) Path route with single path",
			path:       "/api/",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
		},
		{
			title:      "(GET) Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "Handler",
		},
		{
			title:      "(POST) Path route with single path",
			path:       "/api/",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			kind:       "Handler",
		},
		{
			title:      "(DELETE) Path route with single path",
			path:       "/api/",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
			kind:       "Handler",
		},
		{
			title:      "(PUT) Path route with single path",
			path:       "/api/",
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			kind:       "Handler",
		},
		{
			title:      "(Head) Path route with single path",
			path:       "/api/",
			method:     http.MethodHead,
			statusCode: http.StatusOK,
			kind:       "Handler",
		},
		{
			title:      "(Options) Path route with single path",
			path:       "/api/",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "Handler",
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test: %s path: %s method %s", test.title, test.path, test.method), func(t *testing.T) {
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

	switch {
	case rt.method == http.MethodGet && rt.kind == "HandlerFunc":
		fallthrough
	case rt.method == http.MethodPost && rt.kind == "HandlerFunc":
		fallthrough
	case rt.method == http.MethodDelete && rt.kind == "HandlerFunc":
		fallthrough
	case rt.method == http.MethodPut && rt.kind == "HandlerFunc":
		fallthrough
	case rt.method == http.MethodHead && rt.kind == "HandlerFunc":
		fallthrough
	case rt.method == http.MethodOptions && rt.kind == "HandlerFunc":
		r.HandleFunc(rt.method, rt.path, handler)
	case rt.method == http.MethodGet && rt.kind == "Handler":
		fallthrough
	case rt.method == http.MethodPost && rt.kind == "Handler":
		fallthrough
	case rt.method == http.MethodDelete && rt.kind == "Handler":
		fallthrough
	case rt.method == http.MethodPut && rt.kind == "Handler":
		fallthrough
	case rt.method == http.MethodHead && rt.kind == "Handler":
		fallthrough
	case rt.method == http.MethodOptions && rt.kind == "Handler":
		r.Handle(rt.method, rt.path, http.HandlerFunc(handler))
	case rt.method == http.MethodGet:
		r.Get(rt.path, handler)
	case rt.method == http.MethodPost:
		r.Post(rt.path, handler)
	case rt.method == http.MethodDelete:
		r.Delete(rt.path, handler)
	case rt.method == http.MethodPut:
		r.Put(rt.path, handler)
	case rt.method == http.MethodHead:
		r.Head(rt.path, handler)
	case rt.method == http.MethodOptions:
		r.Options(rt.path, handler)
	}

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
