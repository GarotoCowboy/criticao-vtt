package tablehandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/tableDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	tableService "github.com/GarotoCowboy/vttProject/api/service/table"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @BasePath /api/v1

// GetTableHandler
// @Summary Get Table
// @Schemes
// @Description Get a table by ID via query parameter
// @Tags Table
// @Accept json
// @Produce json
// @Param id query int true "Table ID"
// @Success 200 {object} tableDTO.TableResponse "No content"
// @Failure 400 {object} tableDTO.ErrorResponse "Invalid ID supplied"
// @Failure 404 {object} tableDTO.ErrorResponse "User Not Found"
// @Router /table [get]
func GetTableHandler(ctx *gin.Context) {
	idStr := ctx.Query("id")

	if idStr == "" {
		handler.SendError(ctx, http.StatusBadRequest, tableDTO.ErrParamIsRequired("id", "queryParameter").Error())
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		handler.SendError(ctx, http.StatusBadRequest, "id must be a positive integer")
		return
	}

	tableData, err := tableService.GetTable(handler.GetHandlerDB(), uint(id))
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := tableDTO.TableResponse{
		ID:         tableData.ID,
		Name:       tableData.Name,
		Password:   tableData.Password,
		OwnerID:    tableData.OwnerID,
		Owner:      tableData.Owner,
		InviteLink: tableData.InviteLink,
	}

	handler.SendSucess(ctx, "getTable", resp)
}
