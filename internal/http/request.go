package http

import "strings"

type Request struct {
	Method         Method
	Path           string
	Version        string
	Headers        Headers
	Body           []byte
	pathParameters PathParameters
}

func (r *Request) GetPath(path string) string {
	if r.pathParameters.Parameters != nil {
		return r.pathParameters.Parameters[path]
	}
	return ""
}

type Headers struct {
	headers map[string]string
}

func (h *Headers) Get(key string) (string, bool) {
	val, found := h.headers[strings.ToLower(key)]
	return val, found
}

func (h *Headers) Set(key string, value string) {
	if h.headers == nil {
		h.headers = map[string]string{}
	}

	h.headers[strings.ToLower(key)] = value
}
