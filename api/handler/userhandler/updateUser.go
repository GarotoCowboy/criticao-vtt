package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateUserHandler(ctx *gin.Context) {

	request := UpdateUserRequest{}

	ctx.BindJSON(&request)

	if err := request.Validate(); err != nil {
		handler.GetHandlerLogger().ErrorF("Validation error: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	id := ctx.Query("id")
	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "string").Error())
		return
	}

	user := models.User{}
	if err := handler.GetHandlerDB().First(&user, id).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, "user not found")
		return
	}
	//Update user

	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Password != "" {
		user.Password = request.Password
	}
	if request.ImageLink != "" {
		user.ImageLink = request.ImageLink
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Firstname != "" {
		user.Firstname = request.Firstname
	}
	if request.Lastname != "" {
		user.Lastname = request.Lastname
	}
	if err := handler.GetHandlerDB().Save(&user).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("error updating user: %v", err.Error())
		return
	}
	handler.SendSucess(ctx, "update-user", user)
}
