package user

import (
	"errors"
	"github.com/GarotoCowboy/vttProject/api/dto/userDTO"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, req userDTO.CreateUserRequest) (models.User, error) {

	//Validate if request bod is valid
	if err := req.Validate(); err != nil {
		return models.User{}, err
	}

	//E-mail validation
	if err := utils.ValidadeEmail(req.Email); err != nil {
		return models.User{}, err
	}

	var user = models.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
	}

	//hash password to bcrypt
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.User{}, err
	}

	user.Password = string(hashedPassword)

	//Save User
	if err := db.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func DeleteUser(db *gorm.DB, id uint) (models.User, error) {

	//Verify if id is valid
	if id == 0 {
		return models.User{}, errors.New("invalid user ID")
	}

	var user = models.User{}

	//Search first user with inputted id
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err

	}

	//Delete user with inputted id
	if err := db.Delete(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUser(db *gorm.DB, id uint) (models.User, error) {
	if id == 0 {
		return models.User{}, errors.New("invalid user ID")
	}

	var user = models.User{}

	if err := db.Where("id=?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

func ListUsers(db *gorm.DB) ([]models.User, error) {

	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func UpdateUser(db *gorm.DB, id uint, req userDTO.UpdateUserRequest) (models.User, error) {

	if id == 0 {
		return models.User{}, errors.New("invalid user ID")
	}

	userData, err := GetUser(db, id)
	if err != nil {
		return models.User{}, err
	}

	if err := req.Validate(); err != nil {
		return models.User{}, err
	}

	//Update userDTO

	if req.Username != "" {
		userData.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return models.User{}, err
		}

		userData.Password = string(hashedPassword)
		//userData.Password = req.Password
	}
	if req.ImageLink != "" {
		userData.ImageLink = req.ImageLink
	}
	if req.Email != "" {
		if err := utils.ValidadeEmail(req.Email); err != nil {
			return models.User{}, err
		}
		userData.Email = req.Email
	}
	if req.Firstname != "" {
		userData.Firstname = req.Firstname
	}
	if req.Lastname != "" {
		userData.Lastname = req.Lastname
	}
	if err := db.Save(&userData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return userData, nil
}
