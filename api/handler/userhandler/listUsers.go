package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
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
// @Success 200 {object} dto.UserResponse "No content"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /users [get]
func ListUsersHandler(ctx *gin.Context) {

	users := []models.User{}
	if err := handler.GetHandlerDB().Find(&users).Error; err != nil {
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	handler.SendSucess(ctx, "list-users", users)
}
