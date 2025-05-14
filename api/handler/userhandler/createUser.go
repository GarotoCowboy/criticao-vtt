package userhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto"
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/utils"
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
// @Param user body UpdateUserRequest true "User data"
// @Success 200 {object} dto.UserResponse "User Created sucessfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request error"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /user [post]
func CreateUserHandler(ctx *gin.Context) {
	request := dto.CreateUserRequest{}
	//err := ctx.BindJSON(&request)
	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := request.Validate(); err != nil {
		handler.GetHandlerLogger().ErrorF("validation error: %v", err.Error())
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

	//hash user password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error hashing password: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user.Password = string(hashedPassword)

	// generic validation email
	if err := utils.ValidadeEmail(user.Email); err != nil {
		handler.GetHandlerLogger().ErrorF("Error verifying email: %v", err)
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.GetHandlerDB().Create(&user).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating user: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	handler.SendSucess(ctx, "create-user", user)
}
