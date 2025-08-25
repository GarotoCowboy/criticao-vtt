package tableuserhandler

import (
	"net/http"

	"github.com/GarotoCowboy/vttProject/api/dto/tableUserDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	tableUserService "github.com/GarotoCowboy/vttProject/api/service/tableUser"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// CreateTableUserByInviteLinkHandler
// @Summary Create tableUser by inviteLink
// @Schemes
// @Description Create tableUser by inviteLink
// @Tags TableUser
// @Accept json
// @Produce json
// @Param tableUser body tableUserDTO.CreateTableUserInviteLinkRequest true "tableUser data"
// @Success 200 {object} tableUserDTO.TableUserResponse "tableUser Created Sucessfully"
// @Failure 400 {object} tableUserDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} tableUserDTO.ErrorResponse "Internal Server Error"
// @Router /tableUser/inviteLink [post]
func CreateTableUserByInviteLinkHandler(ctx *gin.Context) {
	request := tableUserDTO.CreateTableUserInviteLinkRequest{}

	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	tableUser, err := tableUserService.CreateTableUserByInviteLink(handler.GetHandlerDB(), request)

	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating user: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	handler.SendSucess(ctx, "create-userDTO", tableUser)

	//handler.SendSucess(ctx, "create-tableDTO", fullTable)
}
