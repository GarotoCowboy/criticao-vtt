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

// GetTableUserHandler
// @Summary Get TableUser
// @Schemes
// @Description Get a table tableUser by ID via query parameter
// @Tags TableUser
// @Accept json
// @Produce json
// @Param id query int true "Table User ID"
// @Success 200 {object} tableUserDTO.TableUserResponse "No content"
// @Failure 400 {object} tableUserDTO.ErrorResponse "Invalid ID supplied"
// @Failure 404 {object} tableUserDTO.ErrorResponse "User Not Found"
// @Router /tableUser [get]
func GetTableUserHandler(ctx *gin.Context) {
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

	tableUserData, err := tableUserService.GetTableUser(handler.GetHandlerDB(), uint(id))
	if err != nil {
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := tableUserDTO.TableUserResponse{
		ID:      tableUserData.ID,
		TableID: tableUserData.TableID,
		UserID:  tableUserData.UserID,
		Role:    tableUserData.Role,
		Table:   tableUserData.Table,
		User:    tableUserData.User,
	}

	handler.SendSucess(ctx, "getTable", resp)
}
