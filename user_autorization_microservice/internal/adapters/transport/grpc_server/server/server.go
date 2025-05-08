package server

import (
	"log"
	"net"
	"context"
	. "user_autorization/internal/adapters/transport/grpc_server/generated"

	"google.golang.org/grpc"
)

type core interface {
	LoginUserAndGetToken(userEmail, password string) (string, error)
	RegisterNewUser(userName, userEmail, password string) error
}

type server struct {
	port string
	listener net.Listener
	grpcServer *grpc.Server
	core core
	UserAutorizationServer
}

func NewServer(port string, c core) *server {
	return &server{
		port: port,
		core: c,
	}
}

func (server *server) Start() error {
	// Создание слушателя для порта
	lis, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		log.Printf("Can not open tcp port %v", err)
		return err
	}

	// Инициализация полей и регистрация сервера
	server.listener = lis
	server.grpcServer = grpc.NewServer()
	RegisterUserAutorizationServer(server.grpcServer, server)

	// Старт сервера
	log.Println("Starting gRPC server on :"+server.port)
	go func () {
		server.grpcServer.Serve(server.listener)
	}()
	return nil
}

func (server *server) LoginUserAndGetToken(ctx context.Context, req *LoginUserAndGetTokenRequest) (*LoginUserAndGetTokenResponse, error) {
	log.Printf("email = %v, password = %v", req.UserEmail, req.Password)
	token, err := server.core.LoginUserAndGetToken(req.UserEmail, req.Password)
	return &LoginUserAndGetTokenResponse{ Token: token }, err
}

func (server *server) RegisterNewUser(ctx context.Context, req *RegisterNewUserRequest) (*RegisterNewUserResponse, error) {
	log.Printf("username = %v, password = %v", req.UserName, req.Password)
	err := server.core.RegisterNewUser(req.UserName, req.UserEmail, req.Password)
	return &RegisterNewUserResponse{ Status: "" }, err
}