package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUserHandler(ctx *gin.Context) {
	request := CreateUserRequest{}
	//err := ctx.BindJSON(&request)
	if err := ctx.BindJSON(&request); err != nil {
		handler.GetLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := request.Validate(); err != nil {
		handler.GetLogger().ErrorF("validation error: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Firstname: request.Firstname,
		Lastname:  request.Lastname,
		Email:     request.Email,
		ImageLink: request.ImageLink,
		Password:  request.Password,
		Username:  request.Username,
	}

	if err := handler.GetDB().Create(&user).Error; err != nil {
		handler.GetLogger().ErrorF("Error creating user: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	handler.SendSucess(ctx, "create-user", user)
}
