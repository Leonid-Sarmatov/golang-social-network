package server

import (
	"api_gateway/internal/adapters/grpc_client/user_authorization"
	"api_gateway/internal/adapters/grpc_client/user_follow"
	"api_gateway/internal/adapters/http_server/handlers/login"
	"api_gateway/internal/adapters/http_server/handlers/posts"
	"api_gateway/internal/adapters/http_server/handlers/register"
	"api_gateway/internal/adapters/http_server/handlers/users"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

/* Адаптер для проверки токена авторизации */
type tockenCheker interface {
	ValidateToken(tokenString string) (string, error)
}

/* HTTP-сервер для внешнего API */
type server struct {
	userauthorization.UserAuthorizationGRPC
    userfollow.UserFollowGRPC
	tockenCheker
}

/* Конструктор */
func NewServer(tch tockenCheker, uagrpc userauthorization.UserAuthorizationGRPC, ufgrpc userfollow.UserFollowGRPC) *server {
	return &server{
		tockenCheker: tch,
		UserAuthorizationGRPC: uagrpc,
		UserFollowGRPC: ufgrpc,
	}
}

func (s *server) Init() {
	r := gin.Default()

	// Подключение CORS заголовков
	r.Use(s.CORSMiddleware())

	// Создание группы маршрутов, требующих авторизацию
	authorized := r.Group("/")
	// Если токен просрочен, то перенаправляем на маршрут авторизации
	authorized.Use(s.JWTAuthMiddleware("/login"))

	/* ------ Вызовы перенапрявляемые микросервису авторизации ------ */

	// Аутентификация
	r.POST("/api/login", login.NewLoginHandler(&s.UserAuthorizationGRPC))
	// Регистрация
	r.POST("/api/register", register.NewRegisterHandler(&s.UserAuthorizationGRPC, &s.UserFollowGRPC))

	/* --- Вызовы перенапрявляемые микросервису социального графа --- */

	// Создать пост
	authorized.POST("/api/posts/create", posts.NewAddNewPostHandler(&s.UserFollowGRPC))
	// Получить посты от конкретного пользователя
	authorized.POST("/api/posts/getByUserName", posts.NewGetPostsAddedByUserHandler(&s.UserFollowGRPC))
	// Получить посты для пользователя от всех его подписок
	authorized.POST("/api/posts/intended", posts.NewGetPostsIntendedForTheUserHandler(&s.UserFollowGRPC))
	// Получить всех пользователей
	authorized.POST("/api/users/getAll", users.NewGetAllUsersHandler(&s.UserFollowGRPC))
	// Подписать одного пользователя на другого
	authorized.POST("/api/users/subscribe", users.NewSubscribeUsersHandler(&s.UserFollowGRPC))

	/* -------- Вызовы перенапрявляемые микросервису контента ------- */
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

func (s *server)CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set(
			"Access-Control-Allow-Headers", 
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
		)
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func (s *server) JWTAuthMiddleware(redirectURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Redirect(http.StatusTemporaryRedirect, redirectURL)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		username, err := s.tockenCheker.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v\n", err)
			c.Redirect(http.StatusTemporaryRedirect, redirectURL)
			c.Abort()
			return
		}

		log.Printf("JWT tocken - OK")

		// Передаём данные дальше по цепочке через контекст
		c.Set("username", username)
		c.Next()
	}
}