package main

import (
	"context"
	"fmt"
	"log"
	"net"

	generated "github.com/ViciousKit/course-chat-auth/generated/auth_v1"
	"github.com/ViciousKit/course-chat-auth/internal/config"
	"github.com/ViciousKit/course-chat-auth/storage"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type srv struct {
	generated.UnimplementedAuthV1Server
	Storage *storage.Storage
}

const (
	errorMissingArguments    = "missing arguments"
	errorInternal            = "internal error"
	errorPasswordDoesntMatch = "password doesn't match"
	errorMissingEntity       = "missing requested entity"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("%+v\n", cfg)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Started app at port :%d \n", cfg.GRPC.Port)

	api := &srv{}
	api.Storage = storage.New(cfg.PGUsername, cfg.PGPassword, cfg.PGDatabase, cfg.PGHost, cfg.PGPort)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	generated.RegisterAuthV1Server(grpcServer, api)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func (s *srv) Create(ctx context.Context, req *generated.CreateRequest) (*generated.CreateResponse, error) {
	method := "Create"

	if req.Password == "" || req.Name == "" || req.Email == "" || req.PasswordConfirm == "" {
		return &generated.CreateResponse{}, fmt.Errorf("%s: %s", method, errorMissingArguments)
	}

	if req.Password != req.PasswordConfirm {
		return &generated.CreateResponse{}, fmt.Errorf("%s: %s", method, errorPasswordDoesntMatch)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)

		return &generated.CreateResponse{}, fmt.Errorf("%s: %s", method, errorInternal)
	}

	if err := s.Storage.CreateUser(ctx, req.Name, req.Email, passHash, int(req.Role)); err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf("%s: %s", method, errorInternal)
	}

	return &generated.CreateResponse{}, nil
}

func (s *srv) Get(ctx context.Context, req *generated.GetRequest) (*generated.GetResponse, error) {
	method := "Get"

	user, err := s.Storage.GetUser(ctx, req.Id)
	if err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf("%s: %s", method, errorMissingEntity)
	}

	return &generated.GetResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  generated.UserRole(user.Role),
	}, nil
}

func (s *srv) Update(ctx context.Context, req *generated.UpdateRequest) (*emptypb.Empty, error) {
	method := "Update"

	if req.Name == "" || req.Email == "" || req.Role == 0 || req.Id == 0 {
		return &emptypb.Empty{}, fmt.Errorf("%s: %s", method, errorMissingArguments)
	}

	if err := s.Storage.UpdateUser(ctx, req.Id, req.Name, req.Email, int(req.Role)); err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf("%s: %s", method, errorInternal)
	}

	return &emptypb.Empty{}, nil
}

func (s *srv) Delete(ctx context.Context, req *generated.DeleteRequest) (*emptypb.Empty, error) {
	method := "Delete"

	if req.Id == 0 {
		return &emptypb.Empty{}, fmt.Errorf("%s: %s", method, errorMissingArguments)
	}

	if err := s.Storage.DeleteUser(ctx, req.Id); err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf("%s: %s", method, errorInternal)
	}

	return &emptypb.Empty{}, nil
}
