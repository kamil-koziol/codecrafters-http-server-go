package http

import (
	"io"
	"strings"
)

type Router struct {
	routes []Route
}

type Route struct {
	Method  Method
	Path    string
	Handler HandlerFunc
}

func (r *Route) Matches(path string) (*PathParameters, bool) {
	p := &PathParameters{
		Parameters: map[string]string{},
	}

	routeSplit := strings.Split(r.Path, "/")
	pathSplit := strings.Split(path, "/")

	longerSplit := max(len(routeSplit), len(pathSplit))
	for i := range longerSplit {
		if i == 0 {
			continue
		}

		routeSegment := routeSplit[i]
		if routeSegment == "*" {
			p.Parameters["*"] = strings.Join(pathSplit[i:], "/")
			return p, true
		}

		isVariable := len(routeSegment) >= 2 && routeSegment[0] == '{' && routeSegment[len(routeSplit[i])-1] == '}'
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
	r.routes = append(r.routes, Route{Path: path, Handler: handler, Method: MethodGET})
}

func (r *Router) POST(path string, handler HandlerFunc) {
	r.routes = append(r.routes, Route{Path: path, Handler: handler, Method: MethodPOST})
}

func (r *Router) findRoute(path string, method Method) (*Route, *PathParameters, bool) {
	for _, route := range r.routes {
		if route.Method != method {
			continue
		}
		if params, matches := route.Matches(path); matches {
			return &route, params, true
		}
	}

	return nil, nil, false
}

func (r *Router) Handle(req *Request, w io.Writer) {
	route, pathParameters, found := r.findRoute(req.Path, req.Method)
	if !found {
		_ = WriteResponse(req, w, StatusNotFound, nil, Headers{})
		return
	}

	if pathParameters != nil {
		req.pathParameters = *pathParameters
	}

	route.Handler(req, w)
}
