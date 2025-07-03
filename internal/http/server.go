package http

import (
	"fmt"
	"net"
)

type Server struct{}

func (s *Server) Run(hostport string) error {
	l, err := net.Listen("tcp", hostport)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %w", hostport, err)
	}

	_, err = l.Accept()
	if err != nil {
		return fmt.Errorf("error accepting connection: %w", err)
	}

	return nil
}
