package userfollow


import (
	"log"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "api_gateway/internal/adapters/grpc_client/user_follow/generated"
)

type UserFollowGRPC struct {
	ip string
	port string
	connection *grpc.ClientConn
	userFollowClient UserFollowClient
}

func NewUserAuthorizationClient(ip, port string) *UserFollowGRPC {
	return &UserFollowGRPC{
		ip: ip,
		port: port,
	}
}

func (c *UserFollowGRPC) Start() error {
	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(c.ip+":"+c.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to server: %v", err)
		return err
	}
	c.connection = conn
	c.userFollowClient = NewUserFollowClient(conn)
	log.Printf("Connection: userFollowClient - OK! host = %v, port = %v", c.ip, c.port)
	return nil
}

func (c *UserFollowGRPC) AddNewUser(userName string) (string, error) {
	res, err := c.userFollowClient.AddNewUser(context.Background(), &AddNewUserRequest{
		UserName: userName,
	})
	if err != nil {
		log.Printf("Failed to AddNewUser: %v", err)
		return "", err
	}
	log.Printf("Successfull AddNewUser")
	return res.ResultMessage, nil
}

func (c *UserFollowGRPC) AddNewPost(userName, color string) (string, error) {
	res, err := c.userFollowClient.AddNewPost(context.Background(), &AddNewPostRequest{
		AutorUserName: userName,
		Color: color,
	})
	if err != nil {
		log.Printf("Failed to AddNewPost: %v", err)
		return "", err
	}
	log.Printf("Successfull AddNewPost")
	return res.ResultMessage, nil
}

func (c *UserFollowGRPC) GetPostsAddedByUser(userName string) ([]*Post, error) {
	res, err := c.userFollowClient.GetPostsAddedByUser(context.Background(), &GetPostsAddedByUserRequest{
		UserName: userName,
	})
	if err != nil {
		log.Printf("Failed to GetPostsAddedByUser: %v", err)
		return nil, err
	}
	log.Printf("Successfull GetPostsAddedByUser")
	return res.Posts, nil
}




