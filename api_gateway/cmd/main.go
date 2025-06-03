package main

import (
	"api_gateway/internal/adapters/grpc_client/user_authorization"
	"api_gateway/internal/adapters/grpc_client/user_follow"
	"api_gateway/internal/adapters/http_server/server"
	"api_gateway/internal/adapters/token"
	"os"
	"time"
)

func main() {

	time.Sleep(5 * time.Second)

	host := os.Getenv("USER_AUTORIZATION_HOST")
	port := os.Getenv("USER_AUTORIZATION_PORT")
	host2 := os.Getenv("USER_FOLLOW_HOST")
	port2 := os.Getenv("USER_FOLLOW_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	cli := userauthorization.NewUserAuthorizationClient(host, port)
	cli.Start()

	cli2 := userfollow.NewUserAuthorizationClient(host2, port2)
	cli2.Start()

	tockenCheker := token.NewTokenJWTAdapter(jwtSecret)

	srv := server.NewServer(tockenCheker, *cli, *cli2)
	srv.Init()

	for {

	}

}