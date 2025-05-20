package userhandler

import (
	userDTO "github.com/GarotoCowboy/vttProject/api/dto/userDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	userService "github.com/GarotoCowboy/vttProject/api/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @BasePath /api/v1

// UpdateUserHandler
// @Summary Update user
// @Schemes
// @Description Update User by ID via query parameter
// @Tags User
// @Accept json
// @Produce json
// @Param userDTO body userDTO.UpdateUserRequest true "User data"
// @Param id query int true "User ID"
// @Success 200 {object} userDTO.UserResponse "User Created sucessfully"
// @Failure 400 {object} userDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} userDTO.ErrorResponse "Internal Server Error"
// @Router /user [put]
func UpdateUserHandler(ctx *gin.Context) {

	request := userDTO.UpdateUserRequest{}

	idParam := ctx.Query("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := userService.UpdateUser(handler.GetHandlerDB(), uint(id), request)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error updating user: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	handler.SendSucess(ctx, "Update-user", user)

}
