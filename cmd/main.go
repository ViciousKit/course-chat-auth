package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	storage *storage.Storage
}

const (
	errorMissingArguments    = "missing arguments"
	errorInternal            = "internal error"
	errorPasswordDoesntMatch = "password doesn't match"
	errorMissingEntity       = "missing requested entity"
)

func main() {
	cfg := config.LoadConfig()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Started app at port :%d \n", cfg.GRPC.Port)

	api := &srv{}
	connection := initStorage(cfg.PGUsername, cfg.PGPassword, cfg.PGDatabase, cfg.PGHost, cfg.PGPort)
	api.storage = storage.New(connection)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	generated.RegisterAuthV1Server(grpcServer, api)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func initStorage(user string, password string, dbname string, host string, port int) *pgx.Conn {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Println("Cant connect pg" + err.Error())
		panic(err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		fmt.Println("Cant ping pg" + err.Error())
		panic(err)
	}
	fmt.Println("Connected!")

	return conn
}

func (s *srv) Create(ctx context.Context, req *generated.CreateRequest) (*generated.CreateResponse, error) {
	if req.Password == "" || req.Name == "" || req.Email == "" || req.PasswordConfirm == "" {
		return &generated.CreateResponse{}, fmt.Errorf(errorMissingArguments)
	}

	if req.Password != req.PasswordConfirm {
		return &generated.CreateResponse{}, fmt.Errorf(errorPasswordDoesntMatch)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)

		return &generated.CreateResponse{}, fmt.Errorf(errorInternal)
	}

	id, err := s.storage.CreateUser(ctx, req.Name, req.Email, passHash, int(req.Role))
	if err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf(errorInternal)
	}

	return &generated.CreateResponse{Id: id}, nil
}

func (s *srv) Get(ctx context.Context, req *generated.GetRequest) (*generated.GetResponse, error) {
	user, err := s.storage.GetUser(ctx, req.Id)
	if err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf(errorMissingEntity)
	}

	return &generated.GetResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  generated.UserRole(user.Role),
	}, nil
}

func (s *srv) Update(ctx context.Context, req *generated.UpdateRequest) (*emptypb.Empty, error) {
	if req.Name == "" || req.Email == "" || req.Role == 0 || req.Id == 0 {
		return &emptypb.Empty{}, fmt.Errorf(errorMissingArguments)
	}

	if err := s.storage.UpdateUser(ctx, req.Id, req.Name, req.Email, int(req.Role)); err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf(errorInternal)
	}

	return &emptypb.Empty{}, nil
}

func (s *srv) Delete(ctx context.Context, req *generated.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.storage.DeleteUser(ctx, req.Id); err != nil {
		fmt.Println(err)

		return nil, fmt.Errorf(errorInternal)
	}

	return &emptypb.Empty{}, nil
}
