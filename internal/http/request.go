package http

type Request struct {
	Method         Method
	Path           string
	Version        string
	Headers        *Headers
	Body           []byte
	pathParameters PathParameters
}

func (r *Request) GetPath(path string) string {
	if r.pathParameters.Parameters != nil {
		return r.pathParameters.Parameters[path]
	}
	return ""
}
