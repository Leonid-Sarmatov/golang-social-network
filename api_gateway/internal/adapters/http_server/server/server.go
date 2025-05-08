package server

import (
	"api_gateway/internal/adapters/http_server/handlers/login"
	"api_gateway/internal/adapters/http_server/handlers/register"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

/* HTTP-сервер для внешнего API */
type server struct {
	login.LoginUser
	register.RegisterUser
}

/* Конструктор */
func NewServer(lu login.LoginUser, ru register.RegisterUser) *server {
	return &server{
		LoginUser: lu,
		RegisterUser: ru,
	}
}

func (s *server) Init() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	/* Вызовы перенапрявляемые микросервису авторизации */

	// Аутентификация
	r.POST("/api/login", login.NewLoginHandler(s.LoginUser))
	// Регистрация
	r.POST("/api/register", register.NewLoginHandler(s.RegisterUser))

	/* Вызовы перенапрявляемые микросервису социального графа */
	/* Вызовы перенапрявляемые микросервису контента */
	//r.GET("/", nil)

	// Настройка HTTP-сервера с таймаутами
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,  // Максимальное время чтения запроса
		WriteTimeout: 10 * time.Second, // Максимальное время записи ответа
		IdleTimeout:  15 * time.Second, // Максимальное время ожидания следующего запроса
	}

	// Запуск сервера
	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обработка preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}