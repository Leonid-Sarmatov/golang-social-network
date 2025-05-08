package main

import (
	"api_gateway/internal/adapters/grpc_client/user_authorization"
	"api_gateway/internal/adapters/http_server/server"
	"os"
	"time"
)

func main() {

	time.Sleep(5 * time.Second)

	host := os.Getenv("USER_AUTORIZATION_HOST")
	port := os.Getenv("USER_AUTORIZATION_PORT")

	cli := userauthorization.NewUserAuthorizationClient(host, port)
	cli.Start()

	srv := server.NewServer(cli, cli)
	srv.Init()

	for {

	}

}