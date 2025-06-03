package userauthorization

import (
	"log"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "api_gateway/internal/adapters/grpc_client/user_authorization/generated"
)

type UserAuthorizationGRPC struct {
	ip string
	port string
	connection *grpc.ClientConn
	UserAutorizationClient
}

func NewUserAuthorizationClient(ip, port string) *UserAuthorizationGRPC {
	return &UserAuthorizationGRPC{
		ip: ip,
		port: port,
	}
}

func (c *UserAuthorizationGRPC) Start() error {
	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(c.ip+":"+c.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return err
	}
	c.connection = conn
	c.UserAutorizationClient = NewUserAutorizationClient(conn)
	log.Println("Connection: userAuthorizationClient - OK!")
	return nil
}

func (c *UserAuthorizationGRPC) LoginUser(userEmail, password string) (string, error) {
	res, err := c.LoginUserAndGetToken(context.Background(), &LoginUserAndGetTokenRequest{
		UserEmail: userEmail,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

func (c *UserAuthorizationGRPC) RegisterUser(userName, userEmail, password string) error {
	_, err := c.RegisterNewUser(context.Background(), &RegisterNewUserRequest{
		UserName: userName,
		UserEmail: userEmail,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}




