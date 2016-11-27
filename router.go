package mux

import (
	"net/http"
	"path"
)

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{
		routes: map[string][]*Route{},
		Validatoren: map[string]Validator{
			"method": newMethodValidator(),
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
	routes map[string][]*Route
	// See Router.StrictSlash(). This defines the flag for new routes.
	strictSlash bool
	// See Router.SkipClean(). This defines the flag for new routes.
	skipClean bool
	// see Router.UseEncodedPath(). This defines a flag for all routes.
	useEncodedPath bool
	// see Validator
	Validatoren map[string]Validator
}

// Match matches registered routes against the request.
func (r *Router) triggerMatching(req *http.Request) *Route {

	if routesForMethod, found := r.routes[req.Method]; found {
		for _, route := range routesForMethod {
			if route := route.triggerMatching(req); route != nil {
				return route
			}
		}
	}

	return nil
}

// ServeHTTP dispatches the handler registered in the matched route.
//
// When there is a match, the route variables can be retrieved calling
// mux.Vars(request).
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !r.skipClean {

		path := req.URL.Path

		if r.useEncodedPath {
			path = req.URL.EscapedPath()
		}

		// Clean path to canonical form and redirect.
		if p := cleanPath(path); p != path {
			w.Header().Set("Location", p)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}

	route := r.triggerMatching(req)

	if route == nil {
		r.notFoundHandler().ServeHTTP(w, req)
		return
	}

	req = setCurrentRoute(req, route)

	if route.handler == nil {
		route.handler = r.notFoundHandler()
	}

	route.handler.ServeHTTP(w, req)
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
func (r *Router) NewRoute() *Route {
	route := &Route{
		router:         r,
		strictSlash:    r.strictSlash,
		skipClean:      r.skipClean,
		useEncodedPath: r.useEncodedPath,
	}
	return route
}

// RegisterRoute registers a new route
func (r *Router) RegisterRoute(method string, route *Route) *Route {

	if validator, found := r.Validatoren["method"]; found {
		err := validator.Validate(method)

		if err != nil {
			route.err = newBadRouteError(route, err.Error())
		}
	}

	route.methodName = method

	r.routes[method] = append(r.routes[method], route)
	return route
}

// Handle registers a new route with a matcher for the URL path.
// See Route.Path() and Route.Handler().
func (r *Router) Handle(method string, path string, handler http.Handler) *Route {
	return r.RegisterRoute(method, r.NewRoute().Path(path).Handler(handler))
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
func (r *Router) HandleFunc(method string, path string, HandlerFunc func(http.ResponseWriter, *http.Request)) *Route {
	return r.RegisterRoute(method, r.NewRoute().Path(path).HandlerFunc(HandlerFunc))
}

// Get registers a new get route for the URL path
// See Route.Path() and Route.Handler()
func (r *Router) Get(path string, handlerFunc func(http.ResponseWriter, *http.Request)) *Route {
	return r.RegisterRoute(http.MethodGet, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Put registers a new put route for the URL path
// See Route.Path() and Route.Handler()
func (r *Router) Put(path string, handlerFunc func(http.ResponseWriter, *http.Request)) *Route {
	return r.RegisterRoute(http.MethodPut, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Post registers a new post route for the URL path
// See Route.Path() and Route.Handler()
func (r *Router) Post(path string, handlerFunc func(http.ResponseWriter, *http.Request)) *Route {
	return r.RegisterRoute(http.MethodPost, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Delete registers a new delete route for the URL path
// See Route.Path() and Route.Handler()
func (r *Router) Delete(path string, handlerFunc func(http.ResponseWriter, *http.Request)) *Route {
	return r.RegisterRoute(http.MethodDelete, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Options registers a new options route for the URL path
// See Route.Path() and Route.Handler()
func (r *Router) Options(path string, handlerFunc func(http.ResponseWriter, *http.Request)) *Route {
	return r.RegisterRoute(http.MethodOptions, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// Head registers a new  head route for the URL path
// See Route.Path() and Route.Handler()
func (r *Router) Head(path string, handlerFunc func(http.ResponseWriter, *http.Request)) *Route {
	return r.RegisterRoute(http.MethodHead, r.NewRoute().Path(path).HandlerFunc(handlerFunc))
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
func (r *Router) ListenAndServe(port string) error {
	return http.ListenAndServe(port, r)
}

// HasErrors checks if any errors exists
func (r *Router) HasErrors() (bool, []error) {
	errors := []error{}
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
