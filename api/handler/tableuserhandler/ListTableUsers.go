package tableuserhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/tableUserDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	tableUserService "github.com/GarotoCowboy/vttProject/api/service/tableUser"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// ListTableUsersHandler
// @Summary List TableUsers
// @Schemes
// @Description Get list of tables
// @Tags Table
// @Accept json
// @Produce json
// @Success 200 {object} tableUserDTO.TableUserResponse "No content"
// @Failure 500 {object} tableUserDTO.ErrorResponse "Internal Server Error"
// @Router /tableUser [get]
func ListTableUsersHandler(ctx *gin.Context) {

	tableUsers, err := tableUserService.ListTablesUser(handler.GetHandlerDB())
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	var responses = []tableUserDTO.TableUserResponse{}

	for _, tableU := range tableUsers {
		var resp = tableUserDTO.TableUserResponse{
			ID:      tableU.ID,
			TableID: tableU.TableID,
			UserID:  tableU.UserID,
			Table:   tableU.Table,
			User:    tableU.User,
		}
		responses = append(responses, resp)
	}
	handler.SendSucess(ctx, "list-tables", responses)
}
