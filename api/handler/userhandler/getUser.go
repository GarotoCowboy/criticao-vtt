package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// GetUserHandler
// @Summary Get User
// @Schemes
// @Description Get a user by ID via query parameter
// @Tags User
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {object} dto.UserResponse "No content"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID supplied"
// @Failure 404 {object} dto.ErrorResponse "User Not Found"
// @Router /user [get]
func GetUserHandler(ctx *gin.Context) {
	id := ctx.Query("id")

	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, dto.ErrParamIsRequired("id",
			"queryParameter").Error())
		return
	}

	user := models.User{}

	if err := handler.GetHandlerDB().Where("id=?", id).First(&user).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, err.Error())
	}
	handler.SendSucess(ctx, "get-user", user)
}
