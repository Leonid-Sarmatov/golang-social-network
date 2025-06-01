package main

import (
	"api_gateway/internal/adapters/grpc_client/user_authorization"
	"api_gateway/internal/adapters/http_server/server"
	"api_gateway/internal/adapters/token"
	"os"
	"time"
)

func main() {

	time.Sleep(5 * time.Second)

	host := os.Getenv("USER_AUTORIZATION_HOST")
	port := os.Getenv("USER_AUTORIZATION_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	cli := userauthorization.NewUserAuthorizationClient(host, port)
	cli.Start()

	tockenCheker := token.NewTokenJWTAdapter(jwtSecret)

	srv := server.NewServer(cli, cli, tockenCheker)
	srv.Init()

	for {

	}

}