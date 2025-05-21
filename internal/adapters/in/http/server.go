package http

import (
	"delivery/cmd"
	"delivery/internal/generated/servers"
)

var _ servers.ServerInterface = &Server{}

type Server struct {
	Root cmd.CompositionRoot
}

func NewServer(root cmd.CompositionRoot) (*Server, error) {
	return &Server{root}, nil
}
