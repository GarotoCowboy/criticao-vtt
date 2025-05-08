package userhandler

import (
	"fmt"
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

type UpdateUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	ImageLink string `json:"imagelink"`
}

func (r *UpdateUserRequest) Validate() error {
	//If any field is provided, validation is truthy
	if r.Username != "" || r.Password != "" || r.Email != "" || r.Firstname != "" || r.Lastname != "" ||
		r.ImageLink != "" {
		return nil
	}

	//if none of the fields were provided, return falsy
	return fmt.Errorf("at least one valid field must be provided")
}
