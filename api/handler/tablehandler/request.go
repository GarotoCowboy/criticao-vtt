package tablehandler

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/models"
)

func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type CreateTableRequest struct {
	Name string `json:"name"`

	OwnerID uint `json:"owner_id"`

	Members    []models.TableUser `json:"members"`
	InviteLink string             `json:"invite_link"`
	Password   string             `json:"password"`
}

func (r *CreateTableRequest) Validate() error {

	if r.Name == "" && r.Password == "" && r.OwnerID == 0 {
		return fmt.Errorf("request body is empty")
	}
	if r.Name == "" {
		return errParamIsRequired("name", "string")
	}
	if r.OwnerID == 0 {
		return errParamIsRequired("owner_id", "uint")
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
