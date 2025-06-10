package userfollow

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	. "api_gateway/internal/adapters/grpc_client/user_follow/generated"
)

type UserFollowGRPC struct {
	ip               string
	port             string
	connection       *grpc.ClientConn
	userFollowClient UserFollowClient
}

func NewUserAuthorizationClient(ip, port string) *UserFollowGRPC {
	return &UserFollowGRPC{
		ip:   ip,
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
		Color:         color,
	})
	if err != nil {
		log.Printf("Failed to AddNewPost: %v", err)
		return "", err
	}
	log.Printf("Successfull AddNewPost")
	return res.ResultMessage, nil
}

func (c *UserFollowGRPC) GetPostsAddedByUser(userName string, timeFrom, timeTo time.Time) ([]*Post, error) {
	res, err := c.userFollowClient.GetPostsAddedByUser(context.Background(), &GetPostsAddedByUserRequest{
		UserName: userName,
		TimeFrom: timestamppb.New(timeFrom),
		TimeTo:   timestamppb.New(timeTo),
	})
	if err != nil {
		log.Printf("Failed to GetPostsAddedByUser: %v", err)
		return nil, err
	}
	log.Printf("Successfull GetPostsAddedByUser")
	return res.Posts, nil
}

func (c *UserFollowGRPC) GetAllUsers(requesterUserName string) ([]*User, []bool, error) {
	res, err := c.userFollowClient.GetAllUsers(context.Background(), &GetAllUsersRequest{
		RequesterUserName: requesterUserName,
	})
	if err != nil {
		log.Printf("Failed to GetAllUsers: %v", err)
		return nil, nil, err
	}
	log.Printf("Successfull GetAllUsers")
	return res.Users, res.SubscribeToRequester, nil
}

func (c *UserFollowGRPC) GetPostsIntendedForTheUser(requesterUserName string) ([]*Post, error) {
	res, err := c.userFollowClient.GetPostsIntendedForTheUser(context.Background(), &GetPostsIntendedForTheUserRequest{
		UserName: requesterUserName,
	})
	if err != nil {
		log.Printf("Failed to GetPostsIntendedForTheUser: %v", err)
		return nil, err
	}
	log.Printf("Successfull GetPostsIntendedForTheUser")
	return res.Posts, nil
}

func (c *UserFollowGRPC) SubscribeUsers(userName, subscriberUserName string) error {
	_, err := c.userFollowClient.SubscribeUsers(context.Background(), &SubscribeUsersRequest{
		UserName: userName,
		SubscriberUserName: subscriberUserName,
	})
		if err != nil {
		log.Printf("Failed to SubscribeUsers: %v", err)
		return err
	}
	log.Printf("Successfull SubscribeUsers")
	return nil
}
