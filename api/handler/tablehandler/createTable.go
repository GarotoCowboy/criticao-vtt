package tablehandler

import (
	"net/http"

	tableDTO "github.com/GarotoCowboy/vttProject/api/dto/tableDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	tableService "github.com/GarotoCowboy/vttProject/api/service/table"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// CreateTableHandler
// @Summary Create table
// @Schemes
// @Description Create table
// @Tags Table
// @Accept json
// @Produce json
// @Param table body tableDTO.CreateTableRequest true "table data"
// @Success 200 {object} tableDTO.TableResponse "table Created Sucessfully"
// @Failure 400 {object} tableDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} tableDTO.ErrorResponse "Internal Server Error"
// @Router /table [post]
func CreateTableHandler(ctx *gin.Context) {
	request := tableDTO.CreateTableRequest{}

	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		handler.SendError(ctx, http.StatusUnauthorized, "user_id not found in context")
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		handler.SendError(ctx, http.StatusUnauthorized, "invalid user_id type in context")
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	table, err := tableService.CreateTable(handler.GetHandlerDB(), userID, request)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating table: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	resp := tableDTO.TableResponse{
		ID:         table.ID,
		Name:       table.Name,
		Password:   table.Password,
		OwnerID:    userID,
		InviteLink: table.InviteLink,
		Owner:      table.Owner,
	}

	handler.SendSucess(ctx, "create-table", resp)

	//handler.SendSucess(ctx, "create-tableDTO", fullTable)
}
