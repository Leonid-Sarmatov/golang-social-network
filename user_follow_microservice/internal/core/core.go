package core

import (
	"fmt"
	"time"
)

// Структура пользователя
type User struct {
	ID []byte // Уникальный идентификатор
	UserName string // Имя пользователя
	TimeOfCreate int64 // Время регистрации
	Subscribers []string // Подписчики (слайс с именами пользователей, которые подписаны на пользователя)
	Subscriptions []string // Подписки (слайс с именами пользователей, на которых подписан пользователь)
}

// Структура поста
type Post struct {
	ID []byte // Уникальный идентификатор
	AutorUserName string // Имя пользователя принадлежащее создателю поста
	TimeOfCreate int64 // Время создания
	Color string // Цвет
	LikedThePost []string // Имена пользователей, которые поставили лайк
}

// Абстрактное хранилище постов
type postStorage interface {
	// Добавить пост
	AddNewPost(post *Post) error
	// Поставить посту лайк
	SetPostLike(postID []byte, likedUser string) error
	// Получить количество лайков поста
	GetPostLikes(postID []byte) (int, error)
}

// Абстрактное хранилище пользователей
type userStorage interface {
	// Добавить нового пользователя
	AddNewUser(user *User) error
	// Проверить, существует ли такое имя пользователя в системе или нет
	CheckExistsUserName(userName string) (bool, error)
	// Подписать одного пользователя на другого
	SubscribeUsers(userName, subscriberUserName string) error
}

// Генератор уникальных ID для постов и для пользователей
type idGenerator interface {
	// Сгенерировать и записать ID для поста
	GenAndSetIDForPost(post *Post) error
	// Сгенерировать и записать ID для пользователя
	GenAndSetIDForUser(user *User) error
}

// Ядро приложения, бизнес-логика
type Core struct {
	PostStorage postStorage
	UserStorage userStorage
	IdGenerator idGenerator
}

// Конструктор ядра
func NewCore(ps postStorage, us userStorage, idg idGenerator) *Core {
	return &Core{
		PostStorage: ps,
		UserStorage: us,
		IdGenerator: idg,
	}
}

/*
AddNewUser регистрирует нового пользователя

Аргументы:
  - userName string: Имя пользователя

Возвращает:
  - error: ошибка
*/
func (core *Core)AddNewUser(userName string) error {
	// Проверка уникальности имени пользователя
	exist, err := core.UserStorage.CheckExistsUserName(userName)
	if err != nil {
		return fmt.Errorf("Can not check existence this user name")
	}
	if exist {
		return fmt.Errorf("This user name is already exist")
	}
	// Создание пользователя
	u := &User{
		UserName: userName,
		TimeOfCreate: time.Now().Unix(),
	}
	// Создание ID
	core.IdGenerator.GenAndSetIDForUser(u)
	// Запись в хранилище
	err = core.UserStorage.AddNewUser(u)
	if err != nil {
		return fmt.Errorf("Save user into storage failed")
	}
	return nil
}