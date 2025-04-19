package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/gin-gonic/gin"
)

func CreateUserHandler(ctx *gin.Context) {
	request := struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}{}
	ctx.BindJSON(&request)
	handler.Logger.InfoF("request received %v", request)
}
