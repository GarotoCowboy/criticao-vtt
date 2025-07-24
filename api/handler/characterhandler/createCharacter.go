package characterhandler

import (
	"github.com/GarotoCowboy/vttProject/api/dto/characterDTO"
	"github.com/GarotoCowboy/vttProject/api/handler"
	//"github.com/GarotoCowboy/vttProject/api/service/character"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// CreateCharacterHandler
// @Summary Create a Character
// @Schemes
// @Description Create Character
// @Tags Character
// @Accept json
// @Produce json
// @Param character body characterDTO.CreateCharacterSwaggerRequest true "Character data"
// @Success 200 {object} characterDTO.CharacterResponse "characterhandler Created sucessfully"
// @Failure 400 {object} characterDTO.ErrorResponse "Bad request error"
// @Failure 500 {object} characterDTO.ErrorResponse "Internal Server Error"
// @Router /table/character [post]


func CreateCharacterHandler(ctx *gin.Context) {

	// Initialize an empty struct to hold the incoming JSON request data
	request := characterDTO.CreateCharacterRequest{}

	// Bind the JSON payload from the request body into the 'request' struct
	// If binding fails (invalid JSON or missing fields), log the error and return 400 Bad Request
	if err := ctx.BindJSON(&request); err != nil {
		handler.GetHandlerLogger().ErrorF("Error binding json: %v", err.Error())
		handler.SendError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	// Call the user service to create a new user using the validated request data
	// If the service returns an error, log it and return 500 Internal Server Error
	//persona, err :=  character.CreateCharacter(handler.GetHandlerDB(),request)
	//if err != nil {
	//	handler.GetHandlerLogger().ErrorF("Error creating characterhandler: %v", err.Error())
	//	handler.SendError(ctx, http.StatusInternalServerError, err.Error())
	//	return
	//}

	//handler.SendSucess(ctx, "create-userDTO", persona)

}