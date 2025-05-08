package main

import (
	"log"
	"os"
	"time"
	"user_autorization/internal/adapters/storage"
	"user_autorization/internal/adapters/token"
	"user_autorization/internal/adapters/transport/grpc_server/server"
	"user_autorization/internal/core"
)

func main() {

	time.Sleep(5 * time.Second)

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	//selfHost := os.Getenv("HOST")
	selfPort := os.Getenv("PORT")

	jwtSecret := os.Getenv("JWT_SECRET")

	srg := storage.NewPostgresAdapter()
	err := srg.Start(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Printf("Не удалось инициализировать хранилище: %s", err.Error())
	}

	tg := token.NewTokenJWTAdapter(jwtSecret)

	c := core.NewCore(srg, tg)

	srv := server.NewServer(selfPort, c)
	srv.Start()

	for {

	}
}