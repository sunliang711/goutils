package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Service struct {
	Desc *grpc.ServiceDesc
	Ss   any
}

type Server struct {
	options Options

	server *grpc.Server
	// opts   []grpc.ServerOption
}

func (d *Server) Init(opts ...Option) {
	for _, o := range opts {
		o(&d.options)
	}
}

func New(opts ...Option) *Server {
	options := Options{
		host: "",
		port: 9090,
	}

	srv := &Server{options: options}
	srv.Init(opts...)

	return srv
}

func (srv *Server) Start(services []Service, opts ...grpc.ServerOption) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", srv.options.host, srv.options.port))
	if err != nil {
		return err
	}

	server := grpc.NewServer(opts...)
	for i := range services {
		server.RegisterService(services[i].Desc, services[i].Ss)
	}

	srv.server = server
	err = server.Serve(lis)
	if err != nil {
		return err
	}

	return nil

}

func (srv *Server) Close() {
	srv.server.Stop()
}
