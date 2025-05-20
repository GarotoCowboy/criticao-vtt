package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/userDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	userService "github.com/GarotoCowboy/vttProject/api/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// ListUsersHandler
// @Summary List Users
// @Schemes
// @Description Get list of users
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} userDTO.UserResponse "No content"
// @Failure 500 {object} userDTO.ErrorResponse "Internal Server Error"
// @Router /users [get]
func ListUsersHandler(ctx *gin.Context) {

	users, err := userService.ListUsers(handler.GetHandlerDB())
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	var responses = []userDTO.UserResponse{}

	for _, user := range users {
		var resp = userDTO.UserResponse{
			ID:        user.ID,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Email:     user.Email,
			Username:  user.Username,
			ImageLink: user.ImageLink,
		}
		responses = append(responses, resp)
	}
	handler.SendSucess(ctx, "list-users", responses)
}
