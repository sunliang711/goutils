package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GrpcService struct {
	Desc *grpc.ServiceDesc
	Ss   any
}

type GrpcServer struct {
	port   int
	server *grpc.Server
	opts   []grpc.ServerOption
}

func NewServer(port int, opts ...grpc.ServerOption) *GrpcServer {
	return &GrpcServer{port: port, opts: opts}
}

func (srv *GrpcServer) Start(services []GrpcService) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.port))
	if err != nil {
		return err
	}

	server := grpc.NewServer(srv.opts...)
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

func (srv *GrpcServer) Close() {
	srv.server.Stop()
}
