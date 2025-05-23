package storage

import (
	"database/sql"
	"fmt"
	"user_autorization/internal/core"

	_ "github.com/lib/pq"
)

type postgresAdapter struct {
	connection *sql.DB
}

func NewPostgresAdapter() *postgresAdapter {
	return &postgresAdapter{}
}

func (p *postgresAdapter) Start(host, port, user, password, dbname string) error {
	connectionString := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", 
		host, port, user, password, dbname,
	)
	
	// Пробуем создать соединение с базой данных
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	p.connection = db

	// Создание таблицы пользователей
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users_table (
        id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY, 
        user_name VARCHAR(255),
		user_email VARCHAR(255),
        password VARCHAR(255)
    );`)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresAdapter) IsUserExist(userName string) (bool, error) {
	rows, err := p.connection.Query("SELECT * FROM users_table WHERE user_name=$1", userName)
	if err != nil {
		return false, err
	}

	var u core.User
	var counter int = 0
	for rows.Next() {
		err = rows.Scan(&u.ID, &u.UserName, &u.Password)
		if err != nil {
			return false, err
		}
		counter += 1
	}

	if counter != 0 {
		return true, nil
	}
	return false, nil
}

func (p *postgresAdapter) CreateNewUser(userName, userEmail, password string) error {
	_, err := p.connection.Exec(
		`INSERT INTO users_table (user_name, user_email, password) VALUES ($1, $2, $3)`, 
		userName, userEmail, password,
	)
	return err
}

func (p *postgresAdapter) GetUserPassword(userEmail string) (string, error) {
	rows, err := p.connection.Query("SELECT password FROM users_table WHERE user_email=$1", userEmail)
	if err != nil {
		return "", err
	}

	ps := make([]string, 0)
	for rows.Next() {
		p := ""
		err = rows.Scan(&p)
		if err != nil {
			return "", err
		}
		ps = append(ps, p)
	}

	return ps[0], nil
}

func (p *postgresAdapter) GetUserName(userEmail string) (string, error) {
	rows, err := p.connection.Query("SELECT user_name FROM users_table WHERE user_email=$1", userEmail)
	if err != nil {
		return "", err
	}

	ps := make([]string, 0)
	for rows.Next() {
		p := ""
		err = rows.Scan(&p)
		if err != nil {
			return "", err
		}
		ps = append(ps, p)
	}

	return ps[0], nil
}