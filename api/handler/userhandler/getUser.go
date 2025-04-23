package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserHandler(ctx *gin.Context) {
	id := ctx.Query("id")

	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, errParamIsRequired("id",
			"queryParameter").Error())
		return
	}

	user := models.User{}

	if err := handler.GetDB().Where("id=?", id).First(&user).Error; err != nil {
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
	}
	handler.SendSucess(ctx, "get-user", user)
}
