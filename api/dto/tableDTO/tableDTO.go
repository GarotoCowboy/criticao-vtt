package tableDTO

import (
	"fmt"
	"gorm.io/gorm"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type TableResponse struct {
	ID         uint           `json:"id"`
	Name       string         `json:"firstname"`
	OwnerID    uint           `json:"owner_id"`
	InviteLink string         `json:"invite_link"`
	Owner      interface{}    `json:"owner,omitempty"`
	Password   string         `json:"password"`
	Deleted_At gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// ErrorResponse reports the error in the userDTO request
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type CreateTableRequest struct {
	Name    string `json:"name"`
	OwnerID uint   `json:"owner_id"`
	//InviteLink string             `json:"invite_link"`
	Password string `json:"password"`
}

func (r *CreateTableRequest) Validate() error {

	if r.Name == "" && r.Password == "" && r.OwnerID == 0 {
		return fmt.Errorf("request body is empty")
	}
	if r.Name == "" {
		return ErrParamIsRequired("name", "string")
	}
	if r.OwnerID == 0 {
		return ErrParamIsRequired("owner_id", "uint")
	}
	//if r.InviteLink == ""{
	//	return errParamIsRequired("invite_link", "string")
	//}
	return nil
}

type UpdateTableRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	//Members []models.TableUser `json:"members"`
}

func (r *UpdateTableRequest) Validate() error {
	//If any field is provided, validation is truthy
	if r.Password != "" || r.Name != "" /*|| r.Members != nil*/ {
		return nil
	}

	//if none of the fields were provided, return falsy
	return fmt.Errorf("at least one valid field must be provided")
}
