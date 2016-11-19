package mux

import (
	"net/http"
	"path"
)

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{}
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
	routes []*Route
	// See Router.StrictSlash(). This defines the flag for new routes.
	strictSlash bool
	// See Router.SkipClean(). This defines the flag for new routes.
	skipClean bool
	// see Router.UseEncodedPath(). This defines a flag for all routes.
	useEncodedPath bool
}

// Match matches registered routes against the request.
func (r *Router) triggerMatching(req *http.Request) *Route {

	for _, route := range r.routes {
		if route := route.triggerMatching(req); route != nil {
			return route
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

	if route != nil {
		//req = setVars(req, match.Vars)
		//req = setCurrentRoute(req, match.Route)
	}

	if route.handler == nil && r.NotFoundHandler == nil {
		route.handler = http.NotFoundHandler()
	}

	route.handler.ServeHTTP(w, req)
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
	r.routes = append(r.routes, route)
	return route
}

// Handle registers a new route with a matcher for the URL path.
// See Route.Path() and Route.Handler().
func (r *Router) Handle(path string, handler http.Handler) *Route {
	return r.NewRoute().Path(path).Handler(handler)
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
func (r *Router) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Path(path).HandlerFunc(handler)
}

// Headers registers a new route with a matcher for request header values.
// See Route.Headers().
func (r *Router) Headers(pairs ...string) *Route {
	return r.NewRoute().Headers(pairs...)
}

// MatcherFunc registers a new route with a custom matcher function.
// See Route.MatcherFunc().
func (r *Router) MatcherFunc(f MatcherFunc) *Route {
	return r.NewRoute().MatcherFunc(f)
}

// Methods registers a new route with a matcher for HTTP methods.
// See Route.Methods().
func (r *Router) Methods(methods ...string) *Route {
	return r.NewRoute().Methods(methods...)
}

// Path registers a new route with a matcher for the URL path.
// See Route.Path().
func (r *Router) Path(path string) *Route {
	return r.NewRoute().Path(path)
}

// Schemes registers a new route with a matcher for URL schemes.
// See Route.Schemes().
func (r *Router) Schemes(schemes ...string) *Route {
	return r.NewRoute().Schemes(schemes...)
}

// BuildVarsFunc registers a new route with a custom function for modifying
// route variables before building a URL.
func (r *Router) BuildVarsFunc(f BuildVarsFunc) *Route {
	return r.NewRoute().BuildVarsFunc(f)
}
