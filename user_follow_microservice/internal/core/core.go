package core

import (
	"errors"
	"fmt"
	"time"
	//"log"
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
	ID               []byte   // Уникальный идентификатор
	UserName         string   // Имя пользователя
	TimeOfCreate     int64    // Время регистрации
	Subscribers      []string // Подписчики (слайс с именами пользователей, которые подписаны на пользователя)
	Subscriptions    []string // Подписки (слайс с именами пользователей, на которых подписан пользователь)
	SubscribersNum   int      // Количество подписчиков
	SubscriptionsNum int      // Количество полписок
}

// Структура поста
type Post struct {
	ID            []byte   // Уникальный идентификатор
	AutorUserName string   // Имя пользователя принадлежащее создателю поста
	TimeOfCreate  int64    // Время создания
	Color         string   // Цвет
	LikedThePost  []string // Имена пользователей, которые поставили лайк
}

// Абстрактное хранилище постов
type postStorage interface {
	// Добавить пост
	AddNewPost(post *Post) error
	// Прлучить посты, добавленные определенным пользователем
	GetPostsAddedByUser(username string, timeFrom, timeTo time.Time) ([]*Post, error)
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
func (core *Core) AddNewUser(userName string) error {
	// Проверка входных параметров
	if userName == "" {
		return errors.Join(ErrIncorrectData, fmt.Errorf("пустые строки в качестве аргументов"))
	}
	// Создание пользователя
	u := &User{
		UserName:     userName,
		TimeOfCreate: time.Now().Unix(),
	}
	// Создание ID
	core.IdGenerator.GenAndSetIDForUser(u)
	// Запись в хранилище
	err := core.UserStorage.AddNewUser(u)
	if err != nil {
		return errors.Join(ErrWriteData, err)
	}
	return nil
}

/*
AddNewPost добавляет новый пост

Аргументы:
  - userName string: Имя пользователя
  - color string: Цвет поста

Возвращает:
  - error: ошибка
*/
func (c *Core) AddNewPost(userName, color string) error {
	//log.Printf("<user_follow core.go AddNewPost> name = %v", userName)
	// Проверка входных параметров
	if userName == "" || color == "" {
		return errors.Join(ErrIncorrectData, fmt.Errorf("пустые строки в качестве аргументов"))
	}
	// Заполнение необходимых полей поста
	p := Post{
		AutorUserName: userName,
		TimeOfCreate: time.Now().Unix(),
		Color: color,
		LikedThePost: make([]string, 0),
	}
	// Генерация уникального ID
	err := c.IdGenerator.GenAndSetIDForPost(&p)
	if err != nil {
		return errors.Join(ErrCreateResource, fmt.Errorf("не удалось задать ID при создании поста: %v", err))
	}
	// Отправка на сохранение в хранилище
	err = c.PostStorage.AddNewPost(&p)
	if err != nil {
		return errors.Join(ErrWriteData, fmt.Errorf("не удалось сохранить новый пост: %v", err))
	}

	return nil
}

/*
GetPostsAddedByUser прочитывает посты определенного пользователя

Аргументы:
  - userName string: Имя пользователя
  - color string: Цвет поста

Возвращает:
  - error: ошибка
*/
func (c *Core) GetPostsAddedByUser(userName string, timeFrom, timeTo time.Time) ([]*Post, error) {
	//log.Printf("<user_follow core.go GetPostsAddedByUser> name = %v", userName)
	// Проверка входных параметров
	if userName == "" || !timeTo.After(timeFrom) {
		return nil, errors.Join(ErrIncorrectData, fmt.Errorf("некорректные входные параметры"))
	}
	// Получение постов из хранилища
	posts, err := c.PostStorage.GetPostsAddedByUser(userName, timeFrom, timeTo)
	if err != nil {
		return nil, errors.Join(ErrReadData, fmt.Errorf("невозможно прочитать посты созданные пользователем %v: %v", userName, err))
	}
	return posts, nil
}

func (c *Core) SetPostLike(postID []byte, likedUser string) error {
	return nil
}

func (c *Core) GetPostLikes(postID []byte) (int, error) {
	return -1, nil
}

// type coreInterface interface {
// 	// Добавить пост
// 	AddNewPost(userName, color string) error
// 	// Прлучить посты, добавленные определенным пользователем
// 	GetPostsAddedByUser(userName string) ([]*Post, error)
// 	// Поставить посту лайк
// 	SetPostLike(postID []byte, likedUser string) error
// 	// Получить количество лайков поста
// 	GetPostLikes(postID []byte) (int, error)
// }