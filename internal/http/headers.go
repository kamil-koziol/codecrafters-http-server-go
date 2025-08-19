package http

import (
	"fmt"
	"strings"
)

type Headers struct {
	headers map[string]string
}

func GetDefaultHeaders() *Headers {
	h := NewHeaders()
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/plain")
	h.Set("Content-Length", "0")
	return h
}

func NewHeaders() *Headers {
	return &Headers{headers: map[string]string{}}
}

func (h *Headers) Get(key string) (string, bool) {
	val, found := h.headers[strings.ToLower(key)]
	return val, found
}

func (h *Headers) Set(key string, value string) {
	headerName := strings.ToLower(key)

	if hv, ok := h.Get(key); ok {
		h.headers[headerName] = fmt.Sprintf("%s, %s", value, hv)
	} else {
		h.headers[headerName] = value
	}

}

func (h *Headers) Replace(key string, value string) {
	h.headers[strings.ToLower(key)] = value
}
