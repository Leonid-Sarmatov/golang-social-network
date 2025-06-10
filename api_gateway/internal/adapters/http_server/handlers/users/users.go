package users

import (
	"api_gateway/internal/adapters/grpc_client/user_follow/generated"
	"api_gateway/internal/adapters/http_server/messages"
	"fmt"
	"log"
	"net/http"
	//"time"

	//"google.golang.org/protobuf/encoding/protojson"
	"github.com/gin-gonic/gin"
)

type user struct {
	UserName string `json:"username"`
	SubscribeToRequester bool `json:"subscribe_to_requester"`
}

// type getAllUsersRequest struct {
// 	UserName string `json:"username"`
// }

type getAllUsersResponse struct {
	messages.BaseResponse
	Users []user `json:"users"`
}

type subscribeUsersRequest struct {
	UserName string `json:"username"`
}

type UsersInterface interface {
	GetAllUsers(requesterUserName string) ([]*generated.User, []bool, error)
	SubscribeUsers(userName, subscriberUserName string) error
}

func NewGetAllUsersHandler(users UsersInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Print(" -> NewGetAllUsersHandler")

		un, ok := ctx.Get("username")
		if !ok {
			errString := fmt.Sprintf("Ошибка сервера, не удалось получить список пользователей: %s", "не задано имя пользователя")
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		u, sub, err := users.GetAllUsers(un.(string))
		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось получить список пользователей: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		us := make([]user, len(u))
		for i, val := range u {
			us[i] = user{
				UserName: val.UserName,
				SubscribeToRequester: sub[i],
			} 
			log.Printf("username = %v, bool = %v", val.UserName, sub[i])
		}

		ctx.JSON(http.StatusOK, &getAllUsersResponse{
			BaseResponse: messages.BaseResponse{
				Status: "OK",
			},
			Users: us,
		})
	}
}

func NewSubscribeUsersHandler(users UsersInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Print(" -> NewGetAllUsersHandler")
		req := &subscribeUsersRequest{}

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
			errString := fmt.Sprintf("Ошибка сервера, не удалось подписать пользователей: %s", "не задано имя пользователя")
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		err := users.SubscribeUsers(req.UserName, un.(string))
		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось подписать пользователей: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		ctx.JSON(http.StatusOK, &messages.BaseResponse{
			Status: "OK",
		})
	}
}