package tableuserhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/tableDTO"
	"github.com/GarotoCowboy/vttProject/api/dto/tableUserDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	tableUserService "github.com/GarotoCowboy/vttProject/api/service/tableUser"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @BasePath /api/v1

// DeleteTableUserHandler
// @Summary Delete a TableUser
// @Schemes
// @Description Delete a TableUser by ID via query parameter
// @Tags TableUser
// @Accept json
// @Produce json
// @Param id query int true "TableUser ID"
// @Success 200 {string} string "No content"
// @Failure 400 {object} tableUserDTO.ErrorResponse "Invalid ID supplied"
// @Failure 404 {object} tableUserDTO.ErrorResponse "User Not Found"
// @Failure 500 {object} tableUserDTO.ErrorResponse "Internal Server Error"
// @Router /tableUser [delete]
func DeleteTableUserHandler(ctx *gin.Context) {
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

	tableUserData, err := tableUserService.DeleteTableUser(handler.GetHandlerDB(), uint(id))
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
	}

	resp := tableUserDTO.TableUserResponse{
		ID:      tableUserData.ID,
		TableID: tableUserData.TableID,
		UserID:  tableUserData.UserID,
	}

	handler.SendSucess(ctx, "delete-table", resp)

}
