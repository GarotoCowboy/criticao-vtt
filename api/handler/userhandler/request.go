package userhandler

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/config"
	"gorm.io/gorm"
)

var (
	logger *config.Logger
	db     *gorm.DB
)

func errParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

type CreateUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	ImageLink string `json:"imagelink"`
}

func (r *CreateUserRequest) Validate() error {

	if r.Username == "" && r.Password == "" && r.Email == "" && r.Firstname == "" && r.Lastname == "" {
		return fmt.Errorf("request body is empty")
	}

	//if r == nil{
	//	return  fmt.Errorf("request body is empty")
	//}

	if r.Username == "" {
		return errParamIsRequired("username", "string")
	}

	if r.Password == "" {
		return errParamIsRequired("password", "string")
	}

	if r.Email == "" {
		return errParamIsRequired("email", "string")
	}

	if r.Firstname == "" {
		return errParamIsRequired("firstname", "string")
	}

	if r.Lastname == "" {
		return errParamIsRequired("lastname", "string")
	}

	//if r.ImageLink == ""{
	//	return errParamIsRequired("imageLink", "string")
	//}

	return nil
}
