package mux

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type routeTest struct {
	title      string
	path       string
	method     string
	statusCode int
	kind       string
	queries    map[string][]string
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
			title:      "(GET) Path route with vars",
			path:       "/api/user/2",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/:number", handler)
			},
		},
		{
			title:      "(GET) Path route with vars",
			path:       "/api/user/32",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/:number", handler)
			},
		},
		{
			title:      "(GET) Path route with vars",
			path:       "/api/user/32/article/golang",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/:number/article/:string", handler)
			},
		},
		{
			title:      "(GET) Path route with vars",
			path:       "/api/user/3",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/#([0-9]{1,})", handler)
			},
		},
		{
			title:      "(GET) Path route with vars",
			path:       "/api/artcile?limit=10",
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			queries:    map[string][]string{"limit": []string{"10"}},
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/artcile", handler)
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
		if rt.queries != nil {

			ok := false
			queries := getQueries(r)
		outerLoop:
			for k, v := range rt.queries {
				if reflect.DeepEqual(v, queries.Get(k)) {
					ok = true
					break outerLoop
				}
			}

			if ok {
				w.Write([]byte("unsuccesfully"))
			} else {
				w.Write([]byte("unsuccesfully"))
			}
			return
		}

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

	if res.Code != rt.statusCode || content.String() != "succesfully" {
		return res.Code, false
	}

	return res.Code, true
}

func TestRouteNotfound(t *testing.T) {

	var methods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodPut,
		http.MethodOptions,
		http.MethodConnect,
	}

	for _, method := range methods {
		t.Run(fmt.Sprintf("Method: %s", method), func(t *testing.T) {
			r := NewRouter()
			req, _ := http.NewRequest(method, "http://localhost/echo", nil)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != 404 {
				t.Errorf("Unexpected status code (%d)", res.Code)
			}
		})
	}
}

func TestRouteWithoutHandler(t *testing.T) {

	var methods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodPut,
		http.MethodOptions,
		http.MethodConnect,
	}

	for _, method := range methods {
		t.Run(fmt.Sprintf("Method: %s", method), func(t *testing.T) {
			r := NewRouter()
			r.RegisterRoute(method, r.NewRoute().Path("/echo"))
			req, _ := http.NewRequest(method, "http://localhost/echo", nil)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != 404 {
				t.Errorf("Unexpected status code (%d)", res.Code)
			}
		})
	}
}

func TestHasErrors(t *testing.T) {
	routeA := &Route{
		err: errors.New("Bad route"),
	}
	routeB := &Route{
		err: errors.New("Bad method"),
	}

	r := &Router{}
	r.routes = map[string]routes{}
	r.routes[http.MethodGet] = append(r.routes[http.MethodGet], routeA, routeB)

	if ok, errors := r.HasErrors(); !ok || 0 == len(errors) {
		t.Errorf("Has no errros (Status is %v, How many errors ? %v)", ok, len(errors))
	}
}

func TestSortsRoutes(t *testing.T) {

	kinds := []int{1, 3, 2, 3, 2, 3, 3, 2, 1}

	r := &Router{}
	r.routes = map[string]routes{}

	for _, v := range kinds {
		route := &Route{
			kind: v,
		}

		r.routes[http.MethodGet] = append(r.routes[http.MethodGet], route)
	}

	r.SortRoutes()

	routes := r.routes[http.MethodGet]

	if routes[0].kind != 1 || routes[len(routes)-1].kind != 3 || routes[2].kind != 2 {
		t.Error("Sort of routes is bad")
	}
}
