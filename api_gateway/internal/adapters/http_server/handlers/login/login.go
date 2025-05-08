package login

import (
	"api_gateway/internal/adapters/http_server/messages"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type request struct {
	UserEmail string `json:"email"`
	Password string `json:"password"`
}

type LoginUser interface {
	LoginUser(userEmail, password string) (string, error)
}

func NewLoginHandler(lu LoginUser) gin.HandlerFunc {
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
		tokenString, err := lu.LoginUser(req.UserEmail, req.Password)

		if err != nil {
			errString := fmt.Sprintf("Ошибка сервера, не удалось авторизоваться: %s", err.Error())
			log.Println(errString)
			ctx.JSON(http.StatusBadGateway, &messages.BaseResponse{
				Status:       "Error",
				ErrorMessage: errString,
			})
			return
		}

		ctx.Header("Authorization", "Bearer "+tokenString)
		ctx.JSON(http.StatusOK, &messages.BaseResponse{
			Status: "OK",
		})
	}
}