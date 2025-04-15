package main

import (
	"fmt"
	"context"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()
	uri := "neo4j://neo4j-sn:7687"
	username, password := "neo4j", "bubilda123"

	// Создание драйвера с использованием контекста
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Ошибка создания драйвера: %v", err)
	}
	// Закрытие драйвера с таймаутом
	defer func() {
		closeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := driver.Close(closeCtx); err != nil {
			log.Fatalf("Ошибка закрытия драйвера: %v", err)
		}
	}()

	// Создаем сессию для записи
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: "neo4j",
	})
	defer func() {
		if err = session.Close(ctx); err != nil {
			log.Fatalf("Ошибка закрытия сессии: %v", err)
		}
	}()

	// Выполнение транзакции для создания узла
	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "CREATE (p:Person {name: $name}) RETURN p.name AS name"
		params := map[string]any{"name": "Alice"}

		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		if res.Next(ctx) {
			return res.Record().Values[0], nil
		}
		return nil, res.Err()
	})
	if err != nil {
		log.Fatalf("Ошибка транзакции: %v", err)
	}

	fmt.Printf("Создан узел Person с именем: %v\n", result)
}