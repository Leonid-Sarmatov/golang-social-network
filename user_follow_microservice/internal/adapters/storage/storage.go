package storage

import (
	"context"
	"fmt"
	"log"
	"time"
	"user_follow/internal/core"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// const (
// 	uri          = "neo4j://localhost:7687"//"neo4j://neo4j-sn:7687"
// 	databaseName = "neo4j"
// 	username     = "neo4j"
// 	password     = "bubilda123"
// )

type Neo4jStorage struct {
	dbName string
	driver neo4j.DriverWithContext
}

/*
NewNeo4jStorage конструктор адаптера
к базе данных

Возвращает:
  - Neo4jStorage: структура адаптера
*/
func NewNeo4jStorage() *Neo4jStorage {
	return &Neo4jStorage{}
}

/*
StartConnect запускает процесс
инициализации соединения с БД

Возвращает:
  - error: ошибка
*/
func (neo *Neo4jStorage) StartConnect(host, port, dbName, username, password string) error {
	// Создание драйвера с использованием контекста
	driver, err := neo4j.NewDriverWithContext(
		fmt.Sprintf("neo4j://%s:%s", host, port), 
		neo4j.BasicAuth(username, password, ""),
	)
	if err != nil {
		log.Printf("Ошибка создания драйвера: %v", err)
		return fmt.Errorf("Ошибка создания драйвера: %v", err)
	}

	neo.driver = driver
	neo.dbName = dbName

	return nil
}

/*
CloseConnect корректно закрывает соединения с БД

Возвращает:
  - error: ошибка
*/
func (neo *Neo4jStorage) CloseConnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := neo.driver.Close(ctx); err != nil {
		log.Printf("Ошибка при закрытии драйвера: %v", err)
		return fmt.Errorf("ошибка при закрытии драйвера: %v", err)
	}

	return nil
}

/*
openWriteSession открывает сессию для записи

Аргументы:
  - ctx context.Context: контекст

Возвращает:
  - neo4j.SessionWithContext: интерфейс сессии
*/
func (neo *Neo4jStorage) openWriteSession(ctx context.Context) neo4j.SessionWithContext {
	// Создаем сессию для записи
	session := neo.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: neo.dbName,
	})
	return session
}

/*
openReadSession открывает сессию для стения

Аргументы:
  - ctx context.Context: контекст

Возвращает:
  - neo4j.SessionWithContext: интерфейс сессии
*/
func (neo *Neo4jStorage) openReadSession(ctx context.Context) neo4j.SessionWithContext {
	// Создаем сессию для чтения
	session := neo.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: neo.dbName,
	})
	return session
}

/*
closeSession закрывает сесии на чтение и запись

Аргументы:
  - ctx context.Context: контекст
  - session neo4j.SessionWithContext: сессия

Возвращает:
  - error: ошибка
*/
func (neo *Neo4jStorage) closeSession(ctx context.Context, session neo4j.SessionWithContext) error {
	// Закрываем сессию
	if err := session.Close(ctx); err != nil {
		log.Fatalf("Ошибка закрытия сессии: %v", err)
		return err
	}
	return nil
}

/*
AddNewUser добавить нового пользователя

Аргументы:
  - user *core.User: пользователь

Возвращает:
  - error: ошибка
*/
func (neo *Neo4jStorage) AddNewPost(post *core.Post) error {
	ctx := context.Background()
	s := neo.openWriteSession(ctx)

	defer func() {
		err := neo.closeSession(ctx, s)
		if err != nil {
			log.Printf("Не удалось закрыть сессию после сохранения нового поста: %v", err)
		}
	}()

	// Выполнение транзакции для создания поста
	_, err := s.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "CREATE (p:Post {ID: $id, AutorUserName: $name, TimeOfCreate: $time, Color: $color})"
		params := map[string]any{
			"id":    post.ID,
			"name":  post.AutorUserName,
			"time":  time.Now().Unix(),
			"color": post.Color,
		}

		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return nil, res.Err()
	})
	if err != nil {
		log.Fatalf("Ошибка транзакции: %v", err)
	}

	return nil
}

// Поставить посту лайк
// func (neo *Neo4jStorage)SetPostLike(postID []byte, likedUser string) error {

// }

// Получить количество лайков поста
// func (neo *Neo4jStorage)GetPostLikes(postID []byte) (int, error) {

// }

/*
AddNewUser добавить нового пользователя

Аргументы:
  - user *core.User: пользователь

Возвращает:
  - error: ошибка
*/
func (neo *Neo4jStorage)AddNewUser(user *core.User) error {
	ctx := context.Background()
	s := neo.openWriteSession(ctx)

	defer func() {
		err := neo.closeSession(ctx, s)
		if err != nil {
			log.Printf("Не удалось закрыть сессию после сохранения нового поста: %v", err)
		}
	}()

	// Выполнение транзакции для создания поста
	_, err := s.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "CREATE (u:User {UserName: $name})"
		params := map[string]any{
			"name":  user.UserName,
		}

		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return nil, res.Err()
	})
	if err != nil {
		log.Fatalf("Ошибка транзакции: %v", err)
	}

	return nil
}

// Проверить, существует ли такое имя пользователя в системе или нет
// func (neo *Neo4jStorage)CheckExistsUserName(userName string) (bool, error) {

// }

/*
SubscribeUsers подписывает пользователей

Аргументы:
  - userName string: пользователь, на которого подписываются
  - subscriberUserName string: пользователь подписчик

Возвращает:
  - error: ошибка
*/
func (neo *Neo4jStorage)SubscribeUsers(userName, subscriberUserName string) error {
	ctx := context.Background()
	s := neo.openWriteSession(ctx)

	defer func() {
		err := neo.closeSession(ctx, s)
		if err != nil {
			log.Printf("Не удалось закрыть сессию после сохранения нового поста: %v", err)
		}
	}()

	// Выполнение транзакции для подписки пользователей
	_, err := s.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (a:User), (b:User) WHERE a.UserName = $username1 AND b.UserName = $username2 CREATE (a)-[:SUBSCRIBER]->(b)"
		params := map[string]any{
			"username1": subscriberUserName,
			"username2": userName,
		}

		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		return nil, res.Err()
	})
	if err != nil {
		log.Printf("Ошибка транзакции: %v", err)
		return err
	}

	return nil
}

/*
StartConnect запускает процесс
инициализации соединения с БД

Аргументы:
  - storage BlockchainStorage: абстрактное хранилище
  - hc hashCalulator: абстрактный хэш-калькулятор

Возвращает:
  - *Blockchain: указатель на блокчейн
  - error: ошибка
*/
