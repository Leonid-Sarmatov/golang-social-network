package core

import (
	"errors"
	"fmt"
	//"log"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDatabaseRequest      = errors.New("ошибка запроса к базе данных")
	ErrIncorrectData        = errors.New("некорректные данные")
	ErrReadData             = errors.New("ошибка чтения данных")
	ErrWriteData            = errors.New("ошибка записи данных")
	ErrCreateResource       = errors.New("ошибка создания ресурса")
	ErrModificationResource = errors.New("ошибка изменения ресурса")
)

// Структура пользователя
type User struct {
	ID        int    // Уникальный идентификатор
	UserName  string // Имя пользователя
	UserEmail string // E-mail
	Password  string // Хэшированный пароль
}

// Абстрактное хранилище пользователей
type userStorage interface {
	// Добавить нового пользователя
	CreateNewUser(userName, userEmail, password string) error
	// Получить хэшированный пароль
	GetUserPassword(userEmail string) (string, error)
	// Получить хэшированный пароль
	GetUserName(userEmail string) (string, error)
	// Проверить, существует ли такое имя пользователя в системе или нет
	IsUserExist(userName string) (bool, error)
}

// Генератор токенов
type tokenGenerator interface {
	// Создает токен с полезной нагрузкой и временем действия
	CreateToken(data string, minutes int) (string, error)
}

// Ядро приложения, бизнес-логика
type core struct {
	UserStorage    userStorage
	TokenGenerator tokenGenerator
}

func NewCore(us userStorage, tg tokenGenerator) *core {
	return &core{
		UserStorage:    us,
		TokenGenerator: tg,
	}
}

func (c *core) LoginUserAndGetToken(userEmail, password string) (string, error) {
	hashPassword, err := c.UserStorage.GetUserPassword(userEmail)
	if err != nil {
		return "", ErrReadData
	}

	//log.Printf("Получен хэшированный пароль: %v", hashPassword)

	ok := CompareHashes(hashPassword, password)
	if !ok || password == "" {
		return "", ErrIncorrectData
	}

	//log.Printf("Пароли совпали")

	userName, err := c.UserStorage.GetUserName(userEmail)
	if err != nil {
		return "", ErrReadData
	}

	//log.Printf("Получено имя пользователя: %v", userName)

	t, err := c.TokenGenerator.CreateToken(userName, 15)
	if err != nil {
		return "", ErrCreateResource
	}

	//log.Printf("LoginUserAndGetToken - OK! Токен: %v", t)

	return t, nil
}

func (c *core) RegisterNewUser(userName, userEmail, password string) error {
	exist, err := c.UserStorage.IsUserExist(userName)
	if err != nil {
		return errors.Join(ErrReadData, err)
	}

	if exist || password == "" {
		return errors.Join(ErrIncorrectData, fmt.Errorf("this user name is already exist"))
	}

	hashPassword, err := GenerateHash(password)
	if err != nil {
		return errors.Join(ErrCreateResource, err)
	}

	err = c.UserStorage.CreateNewUser(userName, userEmail, hashPassword)
	if err != nil {
		return errors.Join(ErrWriteData, err)
	}

	return nil
}

/*
GenerateHash  создает хеш из строки
*/
func GenerateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

/*
GenerateHash сравнивает хеш с возможным паролем
*/
func CompareHashes(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
