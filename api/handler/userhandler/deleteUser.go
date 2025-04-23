package userhandler

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteUserHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, errParamIsRequired("id",
			"queryParameter").Error())
		return
	}
	user := models.User{}

	if err := handler.GetDB().First(&user, id).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, fmt.Sprintf("User with id: %s not found", id))
		return
	}
	if err := handler.GetDB().Delete(&user).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, fmt.Sprintf("error deleting user with id: %s", id))
		return
	}
	handler.SendSucess(ctx, "delete-user", user)

}
