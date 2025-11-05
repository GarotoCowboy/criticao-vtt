package tableuserhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/tableUserDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	tableUserService "github.com/GarotoCowboy/vttProject/api/service/tableUser"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// CreateTableUserHandler
// @Summary Create tableUser
// @Schemes
// @Description Create tableUser
// @Tags TableUser
// @Accept json
// @Produce json
// @Param tableUser body tableUserDTO.CreateTableUserRequest true "tableUser data"
// @Success 200 {object} tableUserDTO.TableUserResponse "tableUser Created Sucessfully"
// @Failure 400 {object} tableUserDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} tableUserDTO.ErrorResponse "Internal Server Error"
// @Router /tableUser [post]
func CreateTableUserHandler(ctx *gin.Context) {
	request := tableUserDTO.CreateTableUserRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	tableUser, err := tableUserService.CreateTableUser(handler.GetHandlerDB(), request)

	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating tableUser: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	handler.SendSucess(ctx, "create-userDTO", tableUser)

	//handler.SendSucess(ctx, "create-tableDTO", fullTable)
}
