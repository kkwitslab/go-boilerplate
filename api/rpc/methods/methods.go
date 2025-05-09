package methods

import (
	"context"
	"fmt"

	"github.com/kkwitslab/go-boilerplate/internal/pb"
)

type HelloWorldServer struct {
	pb.UnimplementedHelloWorldServer
}

func (h *HelloWorldServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	name := req.Name
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello, %s From gRPC Server.", name),
	}, nil
}
