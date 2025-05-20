package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/userDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	userService "github.com/GarotoCowboy/vttProject/api/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// CreateUserHandler
// @Summary Create User
// @Schemes
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Param user body userDTO.CreateUserRequest true "User data"
// @Success 200 {object} userDTO.UserResponse "User Created sucessfully"
// @Failure 400 {object} userDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} userDTO.ErrorResponse "Internal Server Error"
// @Router /user [post]
func CreateUserHandler(ctx *gin.Context) {
	request := userDTO.CreateUserRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := userService.CreateUser(handler.GetHandlerDB(), request)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating user: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	handler.SendSucess(ctx, "create-userDTO", user)
}
