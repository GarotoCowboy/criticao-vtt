package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUserHandler(ctx *gin.Context) {
	request := CreateUserRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := request.Validate(); err != nil {
		logger.ErrorF("validation error:%v", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	if err := db.Create(&request).Error; err != nil {
		logger.ErrorF("Error creating user: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}
