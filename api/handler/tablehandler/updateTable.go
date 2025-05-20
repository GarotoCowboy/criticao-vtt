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

// UpdateTableHandler
// @Summary Update table
// @Schemes
// @Description Update Table by ID via query parameter
// @Tags Table
// @Accept json
// @Produce json
// @Param tableDTO body tableDTO.UpdateTableRequest true "Table data"
// @Param id query int true "Table ID"
// @Success 200 {object} tableDTO.TableResponse "Table Created sucessfully"
// @Failure 400 {object} tableDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} tableDTO.ErrorResponse "Internal Server Error"
// @Router /table [put]
func UpdateTableHandler(ctx *gin.Context) {

	request := tableDTO.UpdateTableRequest{}

	idParam := ctx.Query("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, "invalid table ID")
		return
	}

	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := tableService.UpdateTable(handler.GetHandlerDB(), uint(id), request)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error updating table: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	handler.SendSucess(ctx, "Update-table", user)

}
