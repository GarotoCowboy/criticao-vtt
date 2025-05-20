package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/userDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	userService "github.com/GarotoCowboy/vttProject/api/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @BasePath /api/v1

// DeleteUserHandler
// @Summary Delete a User
// @Schemes
// @Description Delete a userDTO by ID via query parameter
// @Tags User
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {string} string "No content"
// @Failure 400 {object} userDTO.ErrorResponse "Invalid ID supplied"
// @Failure 404 {object} userDTO.ErrorResponse "User Not Found"
// @Failure 500 {object} userDTO.ErrorResponse "Internal Server Error"
// @Router /user [delete]
func DeleteUserHandler(ctx *gin.Context) {

	idStr := ctx.Query("id")

	if idStr == "" {
		handler.SendError(ctx, http.StatusBadRequest, userDTO.ErrParamIsRequired("id", "queryParameter").Error())
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		handler.SendError(ctx, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	userData, err := userService.DeleteUser(handler.GetHandlerDB(), uint(id))
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
	}

	resp := userDTO.UserResponse{
		ID:        userData.ID,
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		Email:     userData.Email,
		Username:  userData.Username,
		ImageLink: userData.ImageLink,
	}

	handler.SendSucess(ctx, "delete-user", resp)
}
