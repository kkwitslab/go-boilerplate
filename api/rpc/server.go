package rpc

import (
	"github.com/kkwitslab/go-boilerplate/api/rpc/methods"
	"github.com/kkwitslab/go-boilerplate/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func InitializeGRPCServer() (*grpc.Server, error) {
	grpcServer := grpc.NewServer()

	// register our service methods
	pb.RegisterHelloWorldServer(grpcServer, &methods.HelloWorldServer{})

	// enable server reflection
	reflection.Register(grpcServer)

	return grpcServer, nil
}
