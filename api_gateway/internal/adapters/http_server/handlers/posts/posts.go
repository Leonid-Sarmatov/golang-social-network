package posts

import (
	"api_gateway/internal/adapters/grpc_client/user_follow/generated"
	"api_gateway/internal/adapters/http_server/messages"
	"fmt"
	"log"
	"net/http"
	"time"

	//"google.golang.org/protobuf/encoding/protojson"
	"github.com/gin-gonic/gin"
)

type post struct {
	ID            []byte `json:"id"`
	AutorUserName string `json:"author"`
	TimeOfCreate  int64  `json:"create_at"`
	Color         string `json:"color"`
	Content       string `json:"content"`
}

type addNewPostRequest struct {
	Color string `json:"color"`
}

type getPostsAddedByUserRequest struct {
	TimeFrom time.Time `json:"time_from"`
	TimeTo   time.Time `json:"time_to"`
}

type getPostsAddedByUserResponse struct {
	messages.BaseResponse
	Posts []post `json:"posts"`
}

type PostsInterface interface {
	AddNewPost(userName, color string) (string, error)
	GetPostsAddedByUser(userName string, timeFrom, timeTo time.Time) ([]*generated.Post, error)
	GetPostsIntendedForTheUser(requesterUserName string) ([]*generated.Post, error)
}

func NewAddNewPostHandler(posts PostsInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Print(" -> AddNewPostHandler")
		req := &addNewPostRequest{}

		if err := ctx.ShouldBindJSON(req); err != nil {
			// Если JSON отсутствует или неверный, возвращаем ошибку
			errString := fmt.Sprintf("Некорректный запрос, ошибка JSON парсинга: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadRequest, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		un, ok := ctx.Get("username")
		if !ok {
			errString := fmt.Sprintf("Ошибка сервера, не удалось создать новый пост: %s", "не задано имя пользователя")
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		result, err := posts.AddNewPost(un.(string), req.Color)
		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось создать новый пост: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString + result,
			})
			return
		}

		log.Printf("Новый пост успешно создан!")

		ctx.JSON(http.StatusOK, &messages.BaseResponse{
			Status: "OK",
		})
	}
}

func NewGetPostsAddedByUserHandler(posts PostsInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Print(" -> GetPostsAddedByUserHandler")
		req := &getPostsAddedByUserRequest{}

		if err := ctx.ShouldBindJSON(req); err != nil {
			// Если JSON отсутствует или неверный, возвращаем ошибку
			errString := fmt.Sprintf("Некорректный запрос, ошибка JSON парсинга: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadRequest, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		un, ok := ctx.Get("username")
		if !ok {
			errString := fmt.Sprintf("Ошибка сервера, не удалось получить посты: %s", "не задано имя пользователя")
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		posts, err := posts.GetPostsAddedByUser(un.(string), req.TimeFrom, req.TimeTo)
		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось получить посты: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		response := &getPostsAddedByUserResponse{
			BaseResponse: messages.BaseResponse{
				Status: "OK",
			},
			Posts: make([]post, len(posts)),
		}

		for i, val := range posts {
			response.Posts[i] = post{
				ID:            val.Id,
				AutorUserName: val.AutorUserName,
				TimeOfCreate:  val.TimeOfCreate,
				Color:         val.Color,
			}
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func NewGetPostsIntendedForTheUserHandler(posts PostsInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Print(" -> NewGetPostsIntendedForTheUserHandler")

		un, ok := ctx.Get("username")
		if !ok {
			errString := fmt.Sprintf("Ошибка сервера, не удалось получить посты: %s", "не задано имя пользователя")
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		posts, err := posts.GetPostsIntendedForTheUser(un.(string))
		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось получить посты: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		response := &getPostsAddedByUserResponse{
			BaseResponse: messages.BaseResponse{
				Status: "OK",
			},
			Posts: make([]post, len(posts)),
		}

		for i, val := range posts {
			response.Posts[i] = post{
				ID:            val.Id,
				AutorUserName: val.AutorUserName,
				TimeOfCreate:  val.TimeOfCreate,
				Color:         val.Color,
			}
		}

		ctx.JSON(http.StatusOK, response)
	}
}