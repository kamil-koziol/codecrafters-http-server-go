package http

import (
	"fmt"
	"net"
)

func response(statusCode int, reason string) []byte {
	return fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n\r\n", statusCode, reason)
}

type Server struct{}

func (s *Server) Run(hostport string) error {
	l, err := net.Listen("tcp", hostport)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %w", hostport, err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		return fmt.Errorf("error accepting connection: %w", err)
	}

	_, err = conn.Write(response(200, "OK"))
	if err != nil {
		return fmt.Errorf("unable to write: %w", err)
	}

	return nil
}
