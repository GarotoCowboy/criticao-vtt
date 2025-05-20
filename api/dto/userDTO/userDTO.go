package userDTO

import "fmt"

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

// UserResponse User Response Body
type UserResponse struct {
	ID        uint   `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	ImageLink string `json:"imageLink,omitempty"`
}

// ErrorResponse reports the error in the userDTO request
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// CreateUserRequest Struct to use to create an User
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
		return ErrParamIsRequired("username", "string")
	}

	if r.Password == "" {
		return ErrParamIsRequired("password", "string")
	}

	if r.Email == "" {
		return ErrParamIsRequired("email", "string")
	}

	if r.Firstname == "" {
		return ErrParamIsRequired("firstname", "string")
	}

	if r.Lastname == "" {
		return ErrParamIsRequired("lastname", "string")
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
