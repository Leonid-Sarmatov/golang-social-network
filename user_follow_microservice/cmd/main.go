package main

import (
	//"context" // indirect// indirect// iauthorizationndirect
	//"fmt"
	//"log" // authorization
	//"time"
	//"log"
	"os"
	"time"
	"user_follow/internal/adapters/id_gen"
	"user_follow/internal/adapters/storage"
	"user_follow/internal/adapters/transport/grpc_server/server"
	"user_follow/internal/core"
	//"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {

	time.Sleep(6 * time.Second)

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	db := storage.NewNeo4jStorage()
	db.StartConnect(dbHost, dbPort, dbName, dbUser, dbPassword)
	defer db.CloseConnect()

	// user1 := &core.User{
	// 	UserName: "bubilda",
	// }
	// log.Printf("Сообщение: %v", db.AddNewUser(user1))

	idg := idgen.NewIDGenerator()

	c := core.NewCore(db, db, idg)

	srv := server.NewServer("40001", c)
	srv.Start()


	for {

	}

	/*p := &core.Post{
		AutorUserName: "Бубылда",
		Color: "Коричневый",
	}

	idgen.NewIDGenerator().GenAndSetIDForPost(p)
	log.Printf("Сообщение: %v", db.AddNewPost(p))

	user1 := &core.User{
		UserName: "Петрович",
	}

	user2 := &core.User{
		UserName: "Говночист",
	}

	user3 := &core.User{
		UserName: "Турист",
	}

	log.Printf("Сообщение: %v", db.AddNewUser(user1))
	log.Printf("Сообщение: %v", db.AddNewUser(user3))

	log.Printf("Сообщение: %v", db.SubscribeUsers(user2.UserName, user1.UserName))
	log.Printf("Сообщение: %v", db.SubscribeUsers(user3.UserName, user1.UserName))

	for {

	}*/
}
