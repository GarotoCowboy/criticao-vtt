package characterDTO

import (
	"encoding/json"
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SheetData interface{}

type CharacterResponse struct {
	ID          uint        `json:"id"`
	TableUserID uint        `json:"table_user_id"`
	TableUser   interface{} `json:"tableUser"`
	PlayerName  string      `json:"player_name"`
	Name        string      `json:"character_name"`
	SystemKey   string      `json:"system_key"`
	SheetData   T20SheetDto `json:"sheet_data"`
}

type T20SheetDto struct {
	CharacterID    uint                    `json:"characterID"`
	Attributes     models.Attributes       `json:"attributes"`
	HpPoints       models.HpPoints         `json:"hp_points"`
	ManaPoints     models.ManaPoints       `json:"mana_points"`
	Armor          models.Armor            `json:"armor"`
	CharacterInfo  models.CharacterInfo    `json:"character_info"`
	ClassAndLevel  models.ClassAndLevel    `json:"class_and_level"`
	Attacks        []models.Attack         `json:"attacks" `
	Abilities      []models.Ability        `json:"abilities" `
	Skills         map[string]models.Skill `json:"skills" `
	EquipmentItems []models.EquipmentItem  `json:"equipment_items"`
}

type CreateCharacterSwaggerRequest struct {
	TableUserID uint        `json:"table_user_id" example:"1"`
	PlayerName  string      `json:"player_name" example:"Pedro"`
	Name        string      `json:"character_name" example:"Guts"`
	SystemKey   int         `json:"system_key" example:"1"`
	SheetData   T20SheetDto `json:"sheet_data"`
}
type CreateCharacterRequest struct {
	TableUserID uint             `json:"table_user_id"`
	PlayerName  string           `json:"player_name,omitempty"`
	Name        string           `json:"character_name" gorm:"not null"`
	SystemKey   consts.SystemKey `json:"system_key" gorm:"not null"`
	SheetData   json.RawMessage  `json:"sheet_data" gorm:"type:jsonb"`
}

func (r *CreateCharacterRequest) Validate() error {

	if r.TableUserID == 0 && r.SystemKey == 0 && r.Name == "" {
		return fmt.Errorf("request body is empty")
	}
	if r.TableUserID == 0 {
		return ErrParamIsRequired("tableUserID", "string")
	}
	if r.SystemKey == 0 {
		return ErrParamIsRequired("systemKey", "consts.SystemKey")
	}
	if r.Name == "" {
		return ErrParamIsRequired("name", "string")
	}

	return nil
}

type UpdateCharacterRequest struct {
	TableUserID uint            `json:"table_id"`
	UserID      uint            `json:"user_id"`
	Name        string          `json:"name"`
	PlayerName  consts.Role     `json:"role"`
	SheetData   json.RawMessage `json:"sheet_data"`
}

func (r *UpdateCharacterRequest) Validate() error {
	//If any field is provided, validation is truthy
	if r.TableUserID != 0 || r.UserID != 0 || r.Name != "" || r.PlayerName == 0 || r.SheetData == nil {
		return nil
	}

	//if none of the fields were provided, return falsy
	return fmt.Errorf("at least one valid field must be provided")
}
