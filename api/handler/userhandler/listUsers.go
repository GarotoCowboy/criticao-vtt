package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListUsersHandler(ctx *gin.Context) {

	users := []models.User{}
	if err := handler.GetDB().Find(&users).Error; err != nil {
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	handler.SendSucess(ctx, "list-users", users)
}
