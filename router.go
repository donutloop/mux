package mux

import (
	"net/http"
	"path"
	"sort"
	"strings"
)

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]routes),
		Validatoren: map[string]Validator{
			"method": newMethodValidator(),
			"path":   newPathValidator(),
		},
	}
}

// Router registers routes to be matched and dispatches a handler.
//
// It implements the http.Handler interface, so it can be registered to serve
// requests:
//
//     var router = mux.NewRouter()
//
//     func main() {
//         http.Handle("/", router)
//     }
//
// This will send all incoming requests to the router.
type Router struct {
	// Configurable Handler to be used when no route matches.
	NotFoundHandler http.Handler
	// Routes to be matched, in order.
	routes map[string]routes
	// This defines the flag for new routes.
	StrictSlash bool
	// This defines the flag for new routes.
	SkipClean bool
	// This defines a flag for all routes.
	UseEncodedPath bool
	// see Validator
	Validatoren map[string]Validator
	// This defines a flag for all routes.
	CaseSensitiveURL bool
	// this builds a route
	constructRoute func(*Router) RouteInterface
}

// UseRoute that you can use diffrent instances routes
func (r *Router) UseRoute(constructer func(*Router) RouteInterface) {
	r.constructRoute = constructer
}

// triggerMatching matches registered routes against the request.
func (r *Router) triggerMatching(req *http.Request) RouteInterface {

	if routesForMethod, found := r.routes[req.Method]; found {
		for _, route := range routesForMethod {
			if route := route.Match(req); route != nil {
				return route
			}
		}
	}

	return nil
}

// ServeHTTP dispatches the handler registered in the matched route.
//
// When there is a match, the route variables can be retrieved calling
// mux.GetVars(req).Get(":number") or mux.GetVars(req).GetAll()
//
// and the route queires can be retrieved calling
// mux.GetQueries(req).Get(":number") or mux.GetQueries(req).GetAll()
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.SkipClean {

		path := req.URL.Path

		if r.UseEncodedPath {
			path = req.URL.EscapedPath()
		}

		// Clean path to canonical form and redirect.
		if p := cleanPath(path); p != path {
			w.Header().Set("Location", p)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}

	if !r.CaseSensitiveURL {
		req.URL.Path = strings.ToLower(req.URL.Path)
	}

	route := r.triggerMatching(req)

	if route == nil {
		r.notFoundHandler().ServeHTTP(w, req)
		return
	}

	req = AddCurrentRoute(req, route)
	req = AddQueries(req)

	if route.HasVars() {
		req = AddVars(req, route.ExtractVars(req))
	}

	if !route.HasHandler() {
		route.Handler(r.notFoundHandler())
	}

	route.GetHandler().ServeHTTP(w, req)
}

func (r *Router) notFoundHandler() http.Handler {
	if r.NotFoundHandler == nil {
		return http.NotFoundHandler()
	}

	return r.NotFoundHandler
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
// /net/http/server.go
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}

// NewRoute registers an empty route.
func (r *Router) NewRoute() RouteInterface {
	return r.constructRoute(r)
}

// RegisterRoute registers and validates a new route
func (r *Router) RegisterRoute(method string, route RouteInterface) RouteInterface {

	route.SetMethodName(method)

	for _, validatorKey := range [2]string{"method", "path"} {
		if validator, found := r.Validatoren[validatorKey]; found {

			err := validator.Validate(route)

			if err != nil {
				route.SetError(NewBadRouteError(route, err.Error()))
				break
			}
		}
	}
	r.routes[method] = append(r.routes[method], route)
	return route
}

// Handle registers a new route with a matcher for the URL path.
// See Route.Path() and Route.Handler().
func (r *Router) Handle(method string, path string, handler http.Handler) RouteInterface {
	route := r.NewRoute()
	route.Path(path).Handler(handler)
	return r.RegisterRoute(method, route)
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
func (r *Router) HandleFunc(method string, path string, HandlerFunc func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(method, r.NewRoute().Path(path).HandlerFunc(HandlerFunc))
}

// Get registers a new get route for the URL path
// See Route.Path() and Route.HandlerFunc()
func (r *Router) Get(path string, handlerFunc func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(http.MethodGet, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Put registers a new put route for the URL path
// See Route.Path() and Route.HandlerFunc()
func (r *Router) Put(path string, handlerFunc func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(http.MethodPut, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Post registers a new post route for the URL path
// See Route.Path() and Route.HandlerFunc()
func (r *Router) Post(path string, handlerFunc func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(http.MethodPost, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Delete registers a new delete route for the URL path
// See Route.Path() and Route.HandlerFunc()
func (r *Router) Delete(path string, handlerFunc func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(http.MethodDelete, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Options registers a new options route for the URL path
// See Route.Path() and Route.HandlerFunc()
func (r *Router) Options(path string, handlerFunc func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(http.MethodOptions, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Head registers a new  head route for the URL path
// See Route.Path() and Route.HandlerFunc()
func (r *Router) Head(path string, handlerFunc func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(http.MethodHead, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
func (r *Router) ListenAndServe(port string, callback func(errs []error)) {
	var ok bool
	errs := make([]error, 0)

	if ok, errs = r.HasErrors(); ok {
		callback(errs)
		return
	}

	r.SortRoutes()
	errs = append(errs, http.ListenAndServe(port, r))

	if 0 != len(errs) {
		callback(errs)
	}

	return
}

// HasErrors checks if any errors exists
func (r *Router) HasErrors() (bool, []error) {
	errors := make([]error, 0)
	hasError := false

	for _, v := range r.routes {
		for _, vv := range v {
			if vv.HasError() {
				hasError = true
				errors = append(errors, vv.GetError())
			}
		}
	}

	return hasError, errors
}

// SortRoutes sorts the routes (Rank: RegexPath, PathWithVars, PathNormal)
func (r *Router) SortRoutes() {
	for _, v := range r.routes {
		for _, vv := range v {
			sort.Sort(vv.GetMatchers())
		}
		sort.Sort(v)
	}
}

// routes implements the sort interface (len, swap, less)
// see sort.Sort (Standard Library)
type routes []RouteInterface

func (r routes) Len() int {
	return len(r)
}

func (r routes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r routes) Less(i, j int) bool {
	return r[i].Kind() > r[j].Kind()
}
