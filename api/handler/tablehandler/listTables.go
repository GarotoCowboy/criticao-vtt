package tablehandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/tableDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	tableService "github.com/GarotoCowboy/vttProject/api/service/table"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// ListTablesHandler
// @Summary List Tables
// @Schemes
// @Description Get list of tables
// @Tags Table
// @Accept json
// @Produce json
// @Success 200 {object} tableDTO.TableResponse "No content"
// @Failure 500 {object} tableDTO.TableResponse "Internal Server Error"
// @Router /tables [get]
func ListTablesHandler(ctx *gin.Context) {

	tables, err := tableService.ListTables(handler.GetHandlerDB())
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	var responses = []tableDTO.TableResponse{}

	for _, table := range tables {
		var resp = tableDTO.TableResponse{
			ID:         table.ID,
			InviteLink: table.InviteLink,
			OwnerID:    table.OwnerID,
			Password:   table.Password,
			Name:       table.Name,
		}
		responses = append(responses, resp)
	}
	handler.SendSucess(ctx, "list-tables", responses)
}
