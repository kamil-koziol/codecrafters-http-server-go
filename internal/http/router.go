package http

import (
	"io"
	"strings"
)

type Router struct {
	routes []Route
}

type Route struct {
	Path    string
	Handler HandlerFunc
}

func (r *Route) Matches(path string) (*PathParameters, bool) {
	p := &PathParameters{
		Parameters: map[string]string{},
	}

	routeSplit := strings.Split(r.Path, "/")
	pathSplit := strings.Split(path, "/")

	if len(routeSplit) != len(pathSplit) {
		return nil, false
	}

	for i := range len(routeSplit) {
		if i == 0 {
			continue
		}

		isVariable := len(routeSplit[i]) >= 2 && routeSplit[i][0] == '{' && routeSplit[i][len(routeSplit[i])-1] == '}'

		if isVariable {
			variable := routeSplit[i][1 : len(routeSplit[i])-1]
			p.Parameters[variable] = pathSplit[i]
			continue
		}

		if routeSplit[i] != pathSplit[i] {
			return nil, false
		}
	}

	return p, true
}

type PathParameters struct {
	Parameters map[string]string
}

type HandlerFunc func(*Request, io.Writer)

func (r *Router) GET(path string, handler HandlerFunc) {
	r.routes = append(r.routes, Route{Path: path, Handler: handler})
}

func (r *Router) findRoute(path string) (*Route, *PathParameters, bool) {
	for _, route := range r.routes {
		if params, matches := route.Matches(path); matches {
			return &route, params, true
		}
	}

	return nil, nil, false
}

func (r *Router) Handle(req *Request, w io.Writer) {
	route, pathParameters, found := r.findRoute(req.Path)
	if !found {
		_ = WriteResponse(w, StatusNotFound, nil, Headers{})
		return
	}

	if pathParameters != nil {
		req.pathParameters = *pathParameters
	}

	route.Handler(req, w)
}
