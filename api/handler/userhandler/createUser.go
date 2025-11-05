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
// @Param tableUser body userDTO.CreateUserRequest true "User data"
// @Success 200 {object} userDTO.UserResponse "User Created sucessfully"
// @Failure 400 {object} userDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} userDTO.ErrorResponse "Internal Server Error"
// @Router /tableUser [post]
func CreateUserHandler(ctx *gin.Context) {

	// Initialize an empty struct to hold the incoming JSON request data
	request := userDTO.CreateUserRequest{}

	// Bind the JSON payload from the request body into the 'request' struct
	// If binding fails (invalid JSON or missing fields), log the error and return 400 Bad Request
	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	// Call the tableUser service to create a new tableUser using the validated request data
	// If the service returns an error, log it and return 500 Internal Server Error
	user, err := userService.CreateUser(handler.GetHandlerDB(), request)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating tableUser: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	handler.SendSucess(ctx, "create-userDTO", user)
}
