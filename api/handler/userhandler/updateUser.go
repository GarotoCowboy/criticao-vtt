package userhandler

import (
	"net/http"

	userDTO "github.com/GarotoCowboy/vttProject/api/dto/userDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	userService "github.com/GarotoCowboy/vttProject/api/service/user"
	"github.com/gin-gonic/gin"
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
// @Success 200 {object} userDTO.UserResponse "User Created sucessfully"
// @Failure 400 {object} userDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} userDTO.ErrorResponse "Internal Server Error"
// @Router /user/me [put]
func UpdateUserHandler(ctx *gin.Context) {
	request := userDTO.UpdateUserRequest{}

	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		handler.SendError(ctx, http.StatusBadRequest, "user_id not found in context")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		handler.SendError(ctx, http.StatusBadRequest, "invalid user_id type in context")
		return
	}

	//idParam := ctx.Query("id")
	//id, err := strconv.ParseUint(idParam, 10, 64)
	//if err != nil {
	//	handler.SendError(ctx, http.StatusBadRequest, "invalid user ID")
	//	return
	//}

	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := userService.UpdateUser(handler.GetHandlerDB(), userID, request)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error updating user: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	handler.SendSucess(ctx, "Update-user", user)

}
