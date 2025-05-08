package tablehandler

import (
	"github.com/GarotoCowboy/vttProject/api/handler"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTableHandler(ctx *gin.Context) {
	request := CreateTableRequest{}

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

	link, err := verifyLinkCodeNotExists()
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error to generate link: %v", err.Error())
		handler.SendError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	table := models.Table{
		Name:       request.Name,
		Password:   request.Password,
		OwnerID:    request.OwnerID,
		InviteLink: link,
	}

	hashedPassword, err := utils.HashPassword(table.Password)
	if err != nil {
		handler.GetHandlerLogger().ErrorF("Error hashing password: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	table.Password = string(hashedPassword)

	if err := handler.GetHandlerDB().Create(&table).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating table: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ownerAsTableUser := models.TableUser{
		UserID:  table.OwnerID,
		TableID: table.ID,
		Role:    models.Role(1),
	}

	if err := handler.GetHandlerDB().Create(&ownerAsTableUser).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("Error creating owner tableUser: %v", err.Error())
		// não precisa retornar 500, mas pode logar o erro
	}

	// Recupere a tabela com preload completo
	var fullTable models.Table
	if err := handler.GetHandlerDB().
		Preload("Owner").
		Preload("Members").
		First(&fullTable, table.ID).Error; err != nil {
		handler.GetHandlerLogger().ErrorF("Error preloading table: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	handler.SendSucess(ctx, "create-table", fullTable)
}

func verifyLinkCodeNotExists() (string, error) {
	var count int64
	var link string
	for {
		link = utils.StringWithCharset(10)

		err := handler.GetHandlerDB().
			Model(&models.Table{}).
			Where("invite_link = ?", link).
			Count(&count).Error
		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}
	}
	return link, nil

}
