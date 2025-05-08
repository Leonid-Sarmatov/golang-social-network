package register

import (
	"api_gateway/internal/adapters/http_server/messages"
	"fmt"
	"log"

	//"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type request struct {
	UserName string `json:"username"`
	UserEmail string `json:"email"`
	Password string `json:"password"`
}

type RegisterUser interface {
	RegisterUser(userName, userEmail, password string) error
}

func NewLoginHandler(ru RegisterUser) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &request{}

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
		log.Printf("username = %v, password = %v", req.UserName, req.Password)

		err := ru.RegisterUser(req.UserName, req.UserEmail, req.Password)
		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось зарегистрировать пользователя: %s", err.Error())
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