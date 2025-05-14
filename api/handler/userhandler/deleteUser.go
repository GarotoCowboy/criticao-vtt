package userhandler

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/dto"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// DeleteUserHandler
// @Summary Delete a User
// @Schemes
// @Description Delete a user by ID via query parameter
// @Tags User
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {string} string "No content"
// @Failure 400 {object} dto.ErrorResponse "Invalid ID supplied"
// @Failure 404 {object} dto.ErrorResponse "User Not Found"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /user [delete]
func DeleteUserHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		handler.SendError(ctx, http.StatusBadRequest, dto.ErrParamIsRequired("id",
			"queryParameter").Error())
		return
	}
	user := models.User{}

	if err := handler.GetHandlerDB().First(&user, id).Error; err != nil {
		handler.SendError(ctx, http.StatusNotFound, fmt.Sprintf("User with id: %s not found", id))
		return
	}
	if err := handler.GetHandlerDB().Delete(&user).Error; err != nil {
		handler.SendError(ctx, http.StatusInternalServerError, fmt.Sprintf("error deleting user with id: %s", id))
		return
	}
	handler.SendSucess(ctx, "delete-user", user)

}
