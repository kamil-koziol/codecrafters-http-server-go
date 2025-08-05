package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const CRLF string = "\r\n"

type Status int

const (
	StatusOK       Status = 200
	StatusNotFound Status = 404
)

var statusReasons = map[Status]string{
	StatusOK:       "OK",
	StatusNotFound: "Not Found",
}

func WriteResponse(w io.Writer, status Status, body []byte, headers Headers) error {
	reason := statusReasons[status]
	b := fmt.Appendf(nil, "HTTP/1.1 %d %s%s", status, reason, CRLF)

	if body != nil {
		headers.Set("Content-Length", strconv.Itoa(len(body)))
	}

	if headers.headers != nil {
		for k, v := range headers.headers {
			b = fmt.Appendf(b, "%s: %s%s", k, v, CRLF)
		}
	}

	b = fmt.Append(b, CRLF)
	if body != nil {
		b = fmt.Append(b, string(body))
	}

	_, err := w.Write(b)
	return err
}

type Method string

const (
	MethodGET Method = "GET"
)

func scanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte(CRLF)); i >= 0 {
		return i + len(CRLF), data[0:i], nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

// parseRequest parses HTTP request into [Request]
// Example HTTP request:
// GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\nRequest body
func parseRequest(body io.Reader) (*Request, error) {
	req := Request{}
	scanner := bufio.NewScanner(body)
	scanner.Split(scanCRLF)

	mode := "requestLine"
scanning:
	for scanner.Scan() {
		line := scanner.Text()

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
			if line == "" {
				break scanning
			}

			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid header format: %v", parts)
			}
			headerKey := parts[0]
			headerVal := parts[1][1:len(parts[1])]
			req.Headers.Set(headerKey, headerVal)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
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
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		return
	}

	s.Router.Handle(req, conn)
}
