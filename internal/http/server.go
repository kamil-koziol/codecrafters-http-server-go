package http

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internal/http/encoding"
)

var CRLF = []byte{'\r', '\n'}

type Status int

const (
	StatusOK                  Status = 200
	StatusCreated             Status = 201
	StatusNotFound            Status = 404
	StatusInternalServerError Status = 500
)

var statusReasons = map[Status]string{
	StatusOK:                  "OK",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",
	StatusCreated:             "Created",
}

func appendHeader(b []byte, key string, value string) []byte {
	return fmt.Appendf(b, "%s: %s%s", key, value, string(CRLF))
}

func decoderForEncoding(encodingScheme string) (encoding.Decoder, bool) {
	switch encodingScheme {
	case "gzip":
		return &encoding.GZIPDecoder{}, true
	}

	return nil, false
}

func encoderForEncoding(encodingScheme string) (encoding.Encoder, bool) {
	switch encodingScheme {
	case "gzip":
		return &encoding.GZIPEncoder{}, true
	}

	return nil, false
}

func encoderForEncodings(encodingSchemes []string) (encoding.Encoder, string) {
	for _, encoding := range encodingSchemes {
		if encoder, exists := encoderForEncoding(encoding); exists {
			return encoder, encoding
		}
	}

	return &encoding.PlainEncoder{}, ""
}

func WriteResponse(r *Request, w io.Writer, status Status, body []byte, headers Headers) error {
	reason := statusReasons[status]
	b := fmt.Appendf(nil, "HTTP/1.1 %d %s%s", status, reason, string(CRLF))

	if headers.headers != nil {
		for k, v := range headers.headers {
			b = appendHeader(b, k, v)
		}
	}

	if body != nil {
		acceptEncoding, _ := r.Headers.Get("Accept-Encoding")
		encodings := strings.Split(acceptEncoding, ", ")
		encoder, encoding := encoderForEncodings(encodings)

		encodedBody, err := encoder.Encode(body)
		if err != nil {
			return fmt.Errorf("failed to encode body: %w", err)
		}

		b = appendHeader(b, "Content-Length", strconv.Itoa(len(encodedBody)))
		if encoding != "" {
			b = appendHeader(b, "Content-Encoding", encoding)
		}

		b = append(b, CRLF...)
		b = append(b, encodedBody...)
	} else {
		b = fmt.Append(b, string(CRLF))
	}

	_, err := w.Write(b)
	return err
}

type Method string

const (
	MethodGET  Method = "GET"
	MethodPOST Method = "POST"
)

// parseRequest parses HTTP request into [Request]
// Example HTTP request:
// GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\nRequest body
func parseRequest(body io.Reader) (*Request, error) {
	req := Request{}
	var buf []byte

	mode := "requestLine"
	for {
		ch := make([]byte, 1)
		n, err := body.Read(ch)
		if n == 0 && err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading body: %v", err)
		}
		buf = append(buf, ch...)

		if bytes.HasSuffix(buf, CRLF) && mode != "body" {
			line := string(buf[:len(buf)-len(CRLF)])

		modeSwitch:
			switch mode {
			case "requestLine":
				// Request line
				// GET /index.html HTTP/1.1\r\n
				parts := strings.Split(line, " ")
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid request line format")
				}

				req.Method = Method(parts[0])
				req.Path = parts[1]
				version := parts[2][:len(parts[2])]
				req.Version = version

				mode = "headers"
			case "headers":
				// Headers
				// Host: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n

				if len(line) == 0 {
					mode = "body"
					goto modeSwitch
				}

				parts := strings.SplitN(line, ":", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid header format: %v", parts)
				}
				headerKey := parts[0]
				headerVal := parts[1][1:len(parts[1])]
				req.Headers.Set(headerKey, headerVal)
			case "body":
				if contentLength, exists := req.Headers.Get("Content-Length"); exists {
					length := 0
					fmt.Sscanf(contentLength, "%d", &length)
					bodyContentsBytes := make([]byte, length)
					_, err := io.ReadFull(body, bodyContentsBytes)
					if err != nil {
						return nil, fmt.Errorf("failed to read body: %v", err)
					}

					var bodyDecoder encoding.Decoder = &encoding.PlainDecoder{}
					contentEncoding, exists := req.Headers.Get("Content-Encoding")
					if exists {
						var encodingExists bool
						bodyDecoder, encodingExists = decoderForEncoding(contentEncoding)
						if !encodingExists {
							return nil, fmt.Errorf("unsupported encoding scheme: %s", contentEncoding)
						}
					}

					decoded, err := bodyDecoder.Decode(bodyContentsBytes)
					if err != nil {
						return nil, fmt.Errorf("unable to decode: %w", err)
					}

					req.Body = decoded
				}

				return &req, nil
			}

			buf = nil
		}
	}

	return &req, nil
}

type Server struct {
	Router Router
}

func (s *Server) Run(hostport string) error {
	l, err := net.Listen("tcp", hostport)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %w", hostport, err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		req, err := parseRequest(conn)
		if err != nil {
			return
		}
		s.Router.Handle(req, conn)
	}
}
