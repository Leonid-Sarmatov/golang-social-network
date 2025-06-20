package idgen

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"user_follow/internal/core"
)

type IDGenerator struct {}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (idg *IDGenerator) getHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

/*
GenAndSetIDForPost сериализует пост в байтовое представление

Аргументы:
  - post *core.Post: указатель на пост для которого генерируется ID

Возвращает:
  - error: ошибка
*/
func (idg *IDGenerator) GenAndSetIDForPost(post *core.Post) error {
	postBytes, err := serializePost(post)
	if err != nil {
		return fmt.Errorf("can not serialize this post: %v", err)
	}
	post.ID = idg.getHash(postBytes)
	return nil
}

/*
GenAndSetIDForPost сериализует пост в байтовое представление

Аргументы:
  - post *core.Post: указатель на пост для которого генерируется ID

Возвращает:
  - error: ошибка
*/
func (idg *IDGenerator) GenAndSetIDForUser(user *core.User) error {
	postBytes, err := serializeUser(user)
	if err != nil {
		return fmt.Errorf("can not serialize this user: %v", err)
	}
	user.ID = idg.getHash(postBytes)
	return nil
}

/*
serializationPost сериализует пост в байтовое представление

Аргументы:
  - post *core.Post: указатель на сериализуемый пост

Возвращает:
  - []byte: сериализованный пост
  - error: ошибка
*/
func serializePost(post *core.Post) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Запись имени автора
	if err := writeBytes(buf, []byte(post.AutorUserName)); err != nil {
		return nil, err
	}

	// Запись времени создания
	if err := binary.Write(buf, binary.LittleEndian, post.TimeOfCreate); err != nil {
		return nil, err
	}

	// Запись цвета
	if err := writeBytes(buf, []byte(post.Color)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

/*
serializationPost сериализует пост в байтовое представление

Аргументы:
  - post *core.Post: указатель на сериализуемый пост

Возвращает:
  - []byte: сериализованный пост
  - error: ошибка
*/
func serializeUser(user *core.User) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Запись имени автора
	if err := writeBytes(buf, []byte(user.UserName)); err != nil {
		return nil, err
	}

	// Запись времени регистрации
	if err := binary.Write(buf, binary.LittleEndian, user.TimeOfCreate); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
} 

/*
writeBytes записывает слайс байтов в буфер, 
предварительно записывая его длину (uint32)

Аргументы:
  - *bytes.Buffer: buf указатель на буфер для записи
  - []byte: data данные для записи в буфер

Возвращает:
  - error: ошибка
*/
func writeBytes(buf *bytes.Buffer, data []byte) error {
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	if _, err := buf.Write(data); err != nil {
		return err
	}
	return nil
}