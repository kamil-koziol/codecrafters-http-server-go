package http

import (
	"fmt"
	"net"
)

const CRLF = "\r\n"

type Status int

const (
	StatusOK Status = 200
)

var statusReasons = map[Status]string{
	StatusOK: "OK",
}

func response(status Status) []byte {
	reason := statusReasons[status]
	return fmt.Appendf(nil, "HTTP/1.1 %d %s%s%s", status, reason, CRLF, CRLF)
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

	_, err = conn.Write(response(StatusOK))
	if err != nil {
		return fmt.Errorf("unable to write: %w", err)
	}

	return nil
}
