package tablehandler

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/models"
)

func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type CreateTableRequest struct 
type CreateTableRequest struct {
	Name string `json:"name"`

	OwnerID uint `json:"owner_id"`

	Members []models.TableUser `json:"members"`

	Password   string `json:"password"`
	Password string `json:"password"`
}

	if r.Name == "" && r.Password == "" && r.OwnerID == 0 && r.InviteLink == "" {
	if r.Name == "" && r.Password == "" && r.OwnerID == 0 && r.InviteLink == ""{
		return fmt.Errorf("request body is empty")
	}
	if r.Name == "" {
		return errParamIsRequired("name", "string")
	return errParamIsRequired("name", "string")
	if r.OwnerID == 0 {
	if r.OwnerID == 0{
		return errParamIsRequired("owner_id", "uint")
	}
	//if r.InviteLink == ""{
	//	return errParamIsRequired("invite_link", "string")
	//}

	return nil
}
