package server

import (
	"context"
	"fmt"
	"log"
	"time"
	"net"
	. "user_follow/internal/adapters/transport/grpc_server/generated"
	"user_follow/internal/core"

	"google.golang.org/grpc"
)

type coreInterface interface {
	// Добавить пост
	AddNewPost(userName, color string) error
	// Добавить пользователя
	AddNewUser(userName string) error
	// Прлучить посты, добавленные определенным пользователем
	GetPostsAddedByUser(username string, timeFrom time.Time, timeTo time.Time) ([]*core.Post, error)
	// Получить все посты от всех подписок пользователя
	GetPostsIntendedForTheUser(username string) ([]*core.Post, error)
	// Получить вообще всех пользователей
    GetAllUsers(username string) ([]*core.UserSubscribeToRequesterDecorator, error)
	// Подписать пользователей
	SubscribeUsers(userName, subscriberUserName string) error
	// Поставить посту лайк
	//SetPostLike(postID []byte, likedUser string) error
	// Получить количество лайков поста
	//GetPostLikes(postID []byte) (int, error)
}

type server struct {
	port string
	listener net.Listener
	grpcServer *grpc.Server
	core coreInterface
	UserFollowServer
}

func NewServer(port string, c coreInterface) *server {
	return &server{
		port: port,
		core: c,
	}
}

func (s *server) Start() error {
	// Создание слушателя для порта
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Printf("Can not open tcp port %v", err)
		return err
	}

	// Инициализация полей и регистрация сервера
	s.listener = lis
	s.grpcServer = grpc.NewServer()
	RegisterUserFollowServer(s.grpcServer, s)

	// Старт сервера
	log.Println("<user_follow server.go Start> Starting gRPC server on :"+s.port)
	go func () {
		s.grpcServer.Serve(s.listener)
	}()
	return nil
}


// Создание нового пользователя
func (s *server)AddNewUser(ctx context.Context, req *AddNewUserRequest) (*AddNewUserResponse, error) {
	//log.Printf("<user_follow server.go AddNewUser> name = %v", req.UserName)
	err := s.core.AddNewUser(req.UserName)
	if err != nil {
		return &AddNewUserResponse{ ResultMessage: "ERROR" }, fmt.Errorf("не удалось создать пользователя: %v", err)
	}
	return &AddNewUserResponse{ ResultMessage: "OK" }, nil
}

// Создание нового поста
func (s *server)AddNewPost(ctx context.Context, req *AddNewPostRequest) (*AddNewPostResponse, error) {
	//log.Printf("<user_follow server.go AddNewPost> name = %v, color = %v", req.AutorUserName, req.Color)
	err := s.core.AddNewPost(req.AutorUserName, req.Color)
	if err != nil {
		return &AddNewPostResponse{ ResultMessage: "ERROR" }, fmt.Errorf("не удалось создать пост: %v", err)
	}
	return &AddNewPostResponse{ ResultMessage: "OK" }, nil
}

// Получить созданные пользователем посты
func (s *server)GetPostsAddedByUser(ctx context.Context, req *GetPostsAddedByUserRequest) (*GetPostsAddedByUserResponse, error) {
	//log.Printf("<user_follow server.go GetPostsAddedByUser> name = %v", req.UserName)
	posts, err := s.core.GetPostsAddedByUser(req.UserName, req.TimeFrom.AsTime(), req.TimeTo.AsTime())
	if err != nil {
		return &GetPostsAddedByUserResponse{ Posts: nil }, fmt.Errorf("не удалось получить список постов: %v", err)
	}
	resPosts := make([]*Post, len(posts))
	for i, p := range posts {
		resPosts[i] = &Post{
			Id: p.ID,
			AutorUserName: p.AutorUserName,
			TimeOfCreate: p.TimeOfCreate,
			Color: p.Color,
			//LikedThePost: p.LikedThePost,
		}
	}
	return &GetPostsAddedByUserResponse{ Posts: resPosts }, nil
}

// Получить посты от подписок пользователя
func (s *server)GetPostsIntendedForTheUser(ctx context.Context, req *GetPostsIntendedForTheUserRequest) (*GetPostsIntendedForTheUserResponse, error) {
	posts, err := s.core.GetPostsIntendedForTheUser(req.UserName)
	if err != nil {
		return &GetPostsIntendedForTheUserResponse{ Posts: nil }, fmt.Errorf("не удалось получить список постов: %v", err)
	}
	resPosts := make([]*Post, len(posts))
	for i, p := range posts {
		resPosts[i] = &Post{
			Id: p.ID,
			AutorUserName: p.AutorUserName,
			TimeOfCreate: p.TimeOfCreate,
			Color: p.Color,
			//LikedThePost: p.LikedThePost,
		}
	}
	return &GetPostsIntendedForTheUserResponse{ Posts: resPosts }, nil
}


// Подписать одного пользователя на другого
func (s *server)SubscribeUsers(ctx context.Context, req *SubscribeUsersRequest) (*SubscribeUsersResponse, error) {
	err := s.core.SubscribeUsers(req.UserName, req.SubscriberUserName)
	if err != nil {
		return &SubscribeUsersResponse{ ResultMessage: "ERROR" }, fmt.Errorf("не удалось подписать пользователей: %v", err)
	}
	return &SubscribeUsersResponse{ ResultMessage: "OK" }, nil
}

// Получить всех пользователей
func (s *server)GetAllUsers(ctx context.Context, req *GetAllUsersRequest) (*GetAllUsersResponse, error) {
	users, err := s.core.GetAllUsers(req.RequesterUserName)
	if err != nil {
		return &GetAllUsersResponse{ Users: nil }, fmt.Errorf("не удалось получить список пользователей: %v", err)
	}
	resUsers := make([]*User, len(users))
	subscribe_to_requester := make([]bool, len(users))
	for i, u := range users {
		resUsers[i] = &User{
			UserName: u.UserName,
		}
		subscribe_to_requester[i] = u.SubscribeToRequester
	}
	return &GetAllUsersResponse{ Users: resUsers, SubscribeToRequester: subscribe_to_requester}, nil
}

// // Получить количество подписок и подписчиков пользователя
// func (s *server)GetNumSubscribersAndSubscriptions(ctx context.Context, req *GetSubscribersAndSubscriptionsRequest) (*GetSubscribersAndSubscriptionsResponse, error) {
// 	return nil, nil
// }