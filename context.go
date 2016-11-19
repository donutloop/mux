package mux

import (
	"context"
	"net/http"
)

type contextKey int

const (
	varsKey contextKey = iota
	routeKey
)

// Vars returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	if rv := contextGet(r, varsKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

// CurrentRoute returns the matched route for the current request, if any.
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns, unless the KeepContext option is set on the
// Router.
func CurrentRoute(r *http.Request) *Route {
	if rv := contextGet(r, routeKey); rv != nil {
		return rv.(*Route)
	}
	return nil
}

func setVars(r *http.Request, val interface{}) *http.Request {
	return contextSet(r, varsKey, val)
}

func setCurrentRoute(r *http.Request, val interface{}) *http.Request {
	return contextSet(r, routeKey, val)
}

func contextGet(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func contextSet(r *http.Request, key, val interface{}) *http.Request {
	if val == nil {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), key, val))
}
