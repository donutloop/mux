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
	vars       map[string]string
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
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			vars:       map[string]string{":number": "2"},
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/:number", handler)
			},
		},
		{
			title:      "(GET) Path route with vars",
			path:       "/api/user/32",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			vars:       map[string]string{":number": "32"},
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/:number", handler)
			},
		},
		{
			title:      "(GET) Path route with vars",
			path:       "/api/user/32/article/golang",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			vars:       map[string]string{":number": "32", ":string": "golang"},
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/:number/article/:string", handler)
			},
		},
		{
			title:      "(GET) Path route with vars",
			path:       "/api/user/3",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			kind:       "HandlerFunc",
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.HandleFunc(method, "/api/user/#([0-9]{1,})", handler)
			},
		},
		{
			title:      "(GET) Path route with queries",
			path:       "/api/artcile?limit=10",
			method:     http.MethodGet,
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
			code, message, ok := testRoute(test)

			if !ok {
				t.Errorf("Expected status code %v, Actucal status code %v, Actucal message %v", test.statusCode, code, message)
			}
		})
	}
}

func testRoute(rt routeTest) (int, string, bool) {
	handler := func(w http.ResponseWriter, r *http.Request) {

		okQueries := true
		if nil != rt.queries {
			if !reflect.DeepEqual(rt.queries, GetQueries(r).GetAll()) {
				okQueries = false
			}
		}

		okVars := true
		if nil != rt.vars {
			if !reflect.DeepEqual(rt.vars, GetVars(r).GetAll()) {
				okVars = false
			}
		}

		if route := CurrentRoute(r); route == nil {
			w.Write([]byte(fmt.Sprintf("unsuccesfully (Context route : %v )", route)))
			return
		}

		if nil != rt.queries && nil != rt.vars && okVars && okQueries {
			w.Write([]byte("succesfully"))
		} else if nil != rt.queries && okQueries {
			w.Write([]byte("succesfully"))
		} else if nil != rt.vars && okVars {
			w.Write([]byte("succesfully"))
		} else if nil == rt.queries && nil == rt.vars {
			w.Write([]byte("succesfully"))
		} else {
			w.Write([]byte("unsuccesfully"))
		}
	}

	r := Classic()
	rt.route(r, rt.path, rt.method, handler)

	req, _ := http.NewRequest(rt.method, "http://localhost"+rt.path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var content bytes.Buffer
	_, err := io.Copy(&content, res.Body)

	if err != nil {
		return -1, "", false
	}

	if res.Code != rt.statusCode || content.String() != "succesfully" {
		return res.Code, content.String(), false
	}

	return res.Code, content.String(), true
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
			r := Classic()
			req, _ := http.NewRequest(method, "http://localhost/echo", nil)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != http.StatusNotFound {
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
			r := Classic()
			r.RegisterRoute(method, r.NewRoute().Path("/echo"))
			req, _ := http.NewRequest(method, "http://localhost/echo", nil)
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != http.StatusNotFound {
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

	kinds := []int{0, 2, 1, 2, 1, 2, 2, 1, 0}

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

	if routes[len(routes)-1].Kind() != kindNormalPath || routes[len(routes)-3].Kind() != kindVarsPath || routes[0].Kind() != kindRegexPath {
		t.Errorf("Sort of routes is bad")
	}
}

func TestRouterWithMultiRoutes(t *testing.T) {
	router := Classic()

	handler := func(key string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(key))
		}
	}

	paths := map[string]struct {
		key  string
		path string
	}{
		"/api/user/:number/comment/:number": {
			key:  "0",
			path: "/api/user/1/comment/1",
		},
		"/api/echo": {
			key:  "1",
			path: "/api/echo",
		},
		"/api/user/:number": {
			key:  "2",
			path: "/api/user/1",
		},
		"/api/user/:string": {
			key:  "3",
			path: "/api/user/donutloop",
		},
		"/api/user/#([5-9]{1,1})": {
			key:  "4",
			path: "/api/user/6",
		},
		"/api/article/#([a-z]{1,})": {
			key:  "6",
			path: "/api/article/golang",
		},
		"/api/article/#([0-9]{1,1})": {
			key:  "7",
			path: "/api/article/7",
		},
		"/api/article/#(9[0-9]{1,})": {
			key:  "8",
			path: "/api/article/97",
		},
	}

	for path, pathInfo := range paths {
		router.Get(path, handler(pathInfo.key))
	}

	router.SortRoutes()

	server := httptest.NewServer(router)

	for rawPath, pathInfo := range paths {
		url := server.URL + pathInfo.path
		t.Run(fmt.Sprintf("RawPath: %s, Path: %s Url: %s", rawPath, pathInfo.path, url), func(t *testing.T) {
			res, err := http.Get(url)

			if err != nil {
				t.Errorf("Unexpected error (%s)", err.Error())
			}

			var content bytes.Buffer
			_, err = io.Copy(&content, res.Body)

			if err != nil {
				t.Errorf("Unexpected error (%s)", err.Error())
			}

			defer res.Body.Close()

			if content.String() != pathInfo.key {
				t.Errorf("Path %s: Unexpected path key (Expected path key: %s, Actucal path key %s)", rawPath, pathInfo.key, content.String())
			}
		})
	}
}
