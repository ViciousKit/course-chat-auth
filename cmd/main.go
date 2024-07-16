package main

import (
	"context"
	"fmt"
	"log"
	"net"

	generated "github.com/ViciousKit/course-chat-auth/generated/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type srv struct {
	generated.UnimplementedAuthV1Server
}

const port = 8084

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Started app at port :%d", port)

	api := &srv{}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	generated.RegisterAuthV1Server(grpcServer, api)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func (*srv) Create(ctx context.Context, req *generated.CreateRequest) (*generated.CreateResponse, error) {
	return &generated.CreateResponse{}, nil
}

func (*srv) Get(ctx context.Context, req *generated.GetRequest) (*generated.GetResponse, error) {
	return &generated.GetResponse{}, nil
}

func (*srv) Update(ctx context.Context, req *generated.UpdateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (*srv) Delete(ctx context.Context, req *generated.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil

}
