package tablehandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUserHandler(ctx *gin.Context) {
	request := CreateTableRequest{}
	//err := ctx.BindJSON(&request)
	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := request.Validate(); err != nil {
		handler.GetHandlerLogger().ErrorF("validation error: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	table := models.Table{
		Name:       request.Name,
		Password:   request.Password,
		InviteLink: request.InviteLink,
		OwnerID:    request.OwnerID,
		Members:    nil,
	}

	//Verify if the generated link not exists
	table.InviteLink = utils.StringWithCharset(10)

	//hash user password
	hashedPassword, err := utils.HashPassword(table.Password)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error hashing password: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	table.Password = string(hashedPassword)

	if err := handler.GetHandlerDB().Create(&table).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating user: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	//Vinculing owner with table
	if err := handler.GetHandlerDB().Preload("Owner").First(&table, table.ID).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("Error preloading user: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	handler.SendSucess(ctx, "create-user", table)
}
