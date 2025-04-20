package core

import (

)

// Структура пользователя
type User struct {
	UserName string // Имя пользователя
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

// Ядро приложения, бизнес-логика
type Core struct {
	postStorage
	userStorage
}