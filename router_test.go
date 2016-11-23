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
}

func TestPath(t *testing.T) {

	tests := []routeTest{
		{
			title:      "Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
		{
			title:      "Path route with single path",
			path:       "/api/users/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
		},
		{
			title:      "Path route with single path",
			path:       "/api/echo",
			method:     http.MethodPost,
			statusCode: http.StatusOK,
		},
		{
			title:      "Path route with single path",
			path:       "/api/echo",
			method:     http.MethodDelete,
			statusCode: http.StatusOK,
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

	switch rt.method {
	case http.MethodGet:
		r.Get(rt.path, handler)
	case http.MethodPost:
		r.Post(rt.path, handler)
	case http.MethodDelete:
		r.Delete(rt.path, handler)
	}

	req, _ := http.NewRequest(rt.method, "http://localhost"+rt.path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var content bytes.Buffer
	_, err := io.Copy(&content, res.Body)

	if err != nil {
		return -1, false
	}

	if res.Code != rt.statusCode && content.String() == "succesfully" {
		return res.Code, false
	}

	return res.Code, true
}
