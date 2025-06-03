package server

import (
	"context"
	"fmt"
	"log"
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
	GetPostsAddedByUser(userName string) ([]*core.Post, error)
	// Поставить посту лайк
	SetPostLike(postID []byte, likedUser string) error
	// Получить количество лайков поста
	GetPostLikes(postID []byte) (int, error)
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
	log.Println("Starting gRPC server on :"+s.port)
	go func () {
		s.grpcServer.Serve(s.listener)
	}()
	return nil
}


// Создание нового пользователя
func (s *server)AddNewUser(ctx context.Context, req *AddNewUserRequest) (*AddNewUserResponse, error) {
	log.Printf("<user_follow core.go AddNewUser> name = %v", req.UserName)
	err := s.core.AddNewUser(req.UserName)
	if err != nil {
		return &AddNewUserResponse{ ResultMessage: "ERROR" }, fmt.Errorf("не удалось создать пользователя: %v", err)
	}
	return &AddNewUserResponse{ ResultMessage: "OK" }, nil
}

// Создание нового поста
func (s *server)AddNewPost(ctx context.Context, req *AddNewPostRequest) (*AddNewPostResponse, error) {
	log.Printf("<user_follow core.go AddNewPost> name = %v, color = %v", req.AutorUserName, req.Color)
	err := s.core.AddNewPost(req.AutorUserName, req.Color)
	if err != nil {
		return &AddNewPostResponse{ ResultMessage: "ERROR" }, fmt.Errorf("не удалось создать пост: %v", err)
	}
	return &AddNewPostResponse{ ResultMessage: "OK" }, nil
}

// Получить созданные пользователем посты
func (s *server)GetPostsAddedByUser(ctx context.Context, req *GetPostsAddedByUserRequest) (*GetPostsAddedByUserResponse, error) {
	posts, err := s.core.GetPostsAddedByUser(req.UserName)
	if err != nil {
		return &GetPostsAddedByUserResponse{ Posts: nil }, fmt.Errorf("не удалось Получить список постов: %v", err)
	}
	resPosts := make([]*Post, len(posts))
	for i, p := range posts {
		resPosts[i] = &Post{
			Id: string(p.ID),
			AutorUserName: p.AutorUserName,
			TimeOfCreate: p.TimeOfCreate,
			Color: p.Color,
			LikedThePost: p.LikedThePost,
		}
	}
	return &GetPostsAddedByUserResponse{ Posts: resPosts }, nil
}

// Подписать одного пользователя на другого
func (s *server)SubscribeUsers(ctx context.Context, req *SubscribeUsersRequest) (*SubscribeUsersResponse, error) {
	return nil, nil
}

// Получить количество подписок и подписчиков пользователя
func (s *server)GetNumSubscribersAndSubscriptions(ctx context.Context, req *GetSubscribersAndSubscriptionsRequest) (*GetSubscribersAndSubscriptionsResponse, error) {
	return nil, nil
}