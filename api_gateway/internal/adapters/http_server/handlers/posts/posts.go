package posts

import (
	"api_gateway/internal/adapters/grpc_client/user_follow/generated"
	"api_gateway/internal/adapters/http_server/messages"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addNewPostRequest struct {
	Color    string `json:"color"`
}

type getPostsAddedByUserRequest struct {
	UserName string `json:"username"`
}

type getPostsAddedByUserResponse struct {
	messages.BaseResponse
	Colors []string `json:"colors"`
}

type PostsInterface interface {
	AddNewPost(userName, color string) (string, error)
	GetPostsAddedByUser(userName string) ([]*generated.Post, error)
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

		log.Print("Новый пост успешно создан!")

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

		posts, err := posts.GetPostsAddedByUser(req.UserName)
		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось получить посты пользователя пользователя: %s", err.Error())
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
			Colors: make([]string, len(posts)),
		}

		for i, val := range posts {
			response.Colors[i] = val.Color
		}

		ctx.JSON(http.StatusOK, response)
	}
}
