package tableUserDTO

import (
	"fmt"

	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"gorm.io/gorm"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type TableUserResponse struct {
	ID      uint `json:"id"`
	TableID uint `json:"table_id"`
	UserID  uint `json:"user_id"`
	//Role represents a function that user will be in the table
	//Enum: 1,2
	// 1 to Player, 2 to GameMaster
	Role       consts.Role    `json:"role" binding:"required"`
	User       interface{}    `json:"user"`
	Table      interface{}    `json:"table"`
	Deleted_At gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type CreateTableUserRequest struct {
	TableID uint        `json:"table_id"`
	UserID  uint        `json:"user_id"`
	Role    consts.Role `json:"role"`
}

type CreateTableUserInviteLinkRequest struct {
	InviteLink string      `json:"invite_link"`
	UserID     uint        `json:"user_id"`
	Role       consts.Role `json:"role"`
}

func (r *CreateTableUserInviteLinkRequest) Validate() error {
	if r.InviteLink == "" && r.UserID == 0 && r.Role == 0 {
		return fmt.Errorf("request body is empty")
	}

	if r.InviteLink == "" {
		return ErrParamIsRequired("invite_link", "string")
	}

	if r.UserID == 0 {
		return ErrParamIsRequired("user_id", "uint")
	}

	if r.Role == 0 {
		return ErrParamIsRequired("role", "uint")
	}

	return nil
}

func (r *CreateTableUserRequest) Validate() error {

	if r.UserID == 0 && r.TableID == 0 && r.Role == 0 {
		return fmt.Errorf("request body is empty")
	}
	if r.TableID == 0 {
		return ErrParamIsRequired("table_id", "uint")
	}
	if r.UserID == 0 {
		return ErrParamIsRequired("user_id", "uint")
	}
	if r.Role == 0 {
		return ErrParamIsRequired("role", "uint")
	}

	return nil
}

type UpdateTableUserRequest struct {
	TableID uint        `json:"table_id"`
	UserID  uint        `json:"user_id"`
	Role    consts.Role `json:"role"`
}

func (r *UpdateTableUserRequest) Validate() error {
	//If any field is provided, validation is truthy
	if r.TableID != 0 || r.UserID != 0 || r.Role != 0 {
		return nil
	}

	//if none of the fields were provided, return falsy
	return fmt.Errorf("at least one valid field must be provided")
}
