package utils

import (
	"fmt"
	"regexp"
)

func ValidadeEmail(email string) error {

	//regex to generic email;  example: test@test.test
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !regex.MatchString(email) {
		return fmt.Errorf("email address is not valid")
	}
	return nil
}
