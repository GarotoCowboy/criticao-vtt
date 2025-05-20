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

// DeleteTableHandler
// @Summary Delete a Table
// @Schemes
// @Description Delete a Table by ID via query parameter
// @Tags Table
// @Accept json
// @Produce json
// @Param id query int true "Table ID"
// @Success 200 {string} string "No content"
// @Failure 400 {object} tableDTO.ErrorResponse "Invalid ID supplied"
// @Failure 404 {object} tableDTO.ErrorResponse "User Not Found"
// @Failure 500 {object} tableDTO.ErrorResponse "Internal Server Error"
// @Router /table [delete]
func DeleteTableHandler(ctx *gin.Context) {
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

	tableData, err := tableService.DeleteTable(handler.GetHandlerDB(), uint(id))
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
	}

	resp := tableDTO.TableResponse{
		ID:         tableData.ID,
		Name:       tableData.Name,
		Password:   tableData.Password,
		OwnerID:    tableData.OwnerID,
		InviteLink: tableData.InviteLink,
	}

	handler.SendSucess(ctx, "delete-table", resp)

}
