package userfollow


import (
	"log"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "api_gateway/internal/adapters/grpc_client/user_follow/generated"
)

type userFollowClient struct {
	ip string
	port string
	connection *grpc.ClientConn
	userFollowClient UserFollowClient
}

func NewUserAuthorizationClient(ip, port string) *userFollowClient {
	return &userFollowClient{
		ip: ip,
		port: port,
	}
}

func (c *userFollowClient) Start() error {
	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(c.ip+":"+c.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return err
	}
	c.connection = conn
	c.userFollowClient = NewUserFollowClient(conn)
	log.Printf("Connection: userFollowClient - OK! host = %v, port = %v", c.ip, c.port)
	return nil
}

func (c *userFollowClient) AddNewUser(userName string) (string, error) {
	res, err := c.userFollowClient.AddNewUser(context.Background(), &AddNewUserRequest{
		UserName: userName,
	})
	if err != nil {
		log.Fatalf("Failed to AddNewUser: %v", err)
		return "", err
	}
	log.Fatalf("Successfull AddNewUser")
	return res.ResultMessage, nil
}

func (c *userFollowClient) AddNewPost(userName, color string) (string, error) {
	res, err := c.userFollowClient.AddNewPost(context.Background(), &AddNewPostRequest{
		AutorUserName: userName,
		Color: color,
	})
	if err != nil {
		log.Fatalf("Failed to AddNewPost: %v", err)
		return "", err
	}
	log.Fatalf("Successfull AddNewPost")
	return res.ResultMessage, nil
}

func (c *userFollowClient) GetPostsAddedByUser(userName string) ([]*Post, error) {
	res, err := c.userFollowClient.GetPostsAddedByUser(context.Background(), &GetPostsAddedByUserRequest{
		UserName: userName,
	})
	if err != nil {
		log.Fatalf("Failed to GetPostsAddedByUser: %v", err)
		return nil, err
	}
	log.Fatalf("Successfull GetPostsAddedByUser")
	return res.Posts, nil
}




