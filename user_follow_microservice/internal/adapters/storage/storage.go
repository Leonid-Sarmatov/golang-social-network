package storage

import (
	"context"
	"fmt"
	"log"
	"time"
	"user_follow/internal/core"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	uri          = "neo4j://neo4j-sn:7687"
	databaseName = "neo4j"
	username     = "neo4j"
	password     = "bubilda123"
)

type Neo4jStorage struct {
	ctx    context.Context
	driver neo4j.DriverWithContext
}

/*
NewNeo4jStorage конструктор адаптера
к базе данных

Возвращает:
  - error: ошибка
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
func (neo *Neo4jStorage) StartConnect() error {
	// Создание драйвера с использованием контекста
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Printf("Ошибка создания драйвера: %v", err)
		return fmt.Errorf("Ошибка создания драйвера: %v", err)
	}

	neo.driver = driver

	return nil
}

/*
CloseConnect корректно закрывает
соединения с БД

Возвращает:
  - error: ошибка
*/
func (neo *Neo4jStorage) CloseConnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := neo.driver.Close(ctx); err != nil {
		log.Printf("Ошибка при закрытии драйвера: %v", err)
		return fmt.Errorf("Ошибка при закрытии драйвера: %v", err)
	}

	return nil
}

func (neo *Neo4jStorage) openWriteSession(ctx context.Context) neo4j.SessionWithContext {
	// Создаем сессию для записи
	session := neo.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: databaseName,
	})
	return session
}

func (neo *Neo4jStorage) openReadSession(ctx context.Context) neo4j.SessionWithContext {
	// Создаем сессию для чтения
	session := neo.driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeRead,
		DatabaseName: databaseName,
	})
	return session
}

func (neo *Neo4jStorage) closeSession(ctx context.Context, session neo4j.SessionWithContext) error {
	// Закрываем сессию
	if err := session.Close(ctx); err != nil {
		log.Fatalf("Ошибка закрытия сессии: %v", err)
	}
	return nil
}

// Добавить пост
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
		query := "CREATE (p:Post {ID: $id, AutorUserName: $name, TimeOfCreate: $time, Color: $color}) RETURN p"
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

// Добавить нового пользователя
// func (neo *Neo4jStorage)AddNewUser(user *core.User) error {

// }

// Проверить, существует ли такое имя пользователя в системе или нет
// func (neo *Neo4jStorage)CheckExistsUserName(userName string) (bool, error) {

// }

// Подписать одного пользователя на другого
// func (neo *Neo4jStorage)SubscribeUsers(userName, subscriberUserName string) error {

// }

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
