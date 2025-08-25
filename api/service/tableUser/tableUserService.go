package tableUser

import (
	"errors"

	"github.com/GarotoCowboy/vttProject/api/dto/tableUserDTO"
	"github.com/GarotoCowboy/vttProject/api/models"
	"gorm.io/gorm"
)

func CreateTableUser(db *gorm.DB, req tableUserDTO.CreateTableUserRequest) (models.TableUser, error) {

	//Validate if request bod is valid
	if err := req.Validate(); err != nil {
		return models.TableUser{}, err
	}

	var tableUser = models.TableUser{
		UserID:  req.UserID,
		TableID: req.TableID,
		Role:    req.Role,
	}

	//Save User
	if err := db.Create(&tableUser).Error; err != nil {
		return models.TableUser{}, err
	}
	return tableUser, nil
}

func CreateTableUserByInviteLink(db *gorm.DB, req tableUserDTO.CreateTableUserInviteLinkRequest) (models.TableUser, error) {
	if err := req.Validate(); err != nil {
		return models.TableUser{}, err
	}

	var table = models.Table{}
	if err := db.Where("invite_link = ?", req.InviteLink).First(&table).Error; err != nil {
		return models.TableUser{}, err
	}

	var tableUser = models.TableUser{
		UserID:  req.UserID,
		TableID: table.ID,
		Role:    req.Role,
	}
	if err := db.Create(&tableUser).Error; err != nil {
		return models.TableUser{}, err
	}
	return tableUser, nil
}

func DeleteTableUser(db *gorm.DB, id uint) (models.TableUser, error) {

	//Verify if id is valid
	if id == 0 {
		return models.TableUser{}, errors.New("invalid tableUser ID")
	}

	var tableUser = models.TableUser{}

	//Search first table with inputted id
	if err := db.First(&tableUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TableUser{}, errors.New("tableUser not found")
		}
		return models.TableUser{}, err

	}

	//Delete table with inputted id
	if err := db.Delete(&tableUser).Error; err != nil {
		return models.TableUser{}, err
	}
	return tableUser, nil
}

func GetTableUser(db *gorm.DB, id uint) (models.TableUser, error) {
	if id == 0 {
		return models.TableUser{}, errors.New("invalid table_user_id")
	}

	var tableUser = models.TableUser{}

	if err := db.Preload("User").
		Preload("Table").
		Where("id=?", id).First(&tableUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TableUser{}, errors.New("table_user not found")
		}
		return models.TableUser{}, err
	}

	return tableUser, nil
}

func ListTablesUser(db *gorm.DB) ([]models.TableUser, error) {

	var tableUsers []models.TableUser

	if err := db.Find(&tableUsers).Error; err != nil {
		return tableUsers, err
	}
	return tableUsers, nil
}

func UpdateTableUser(db *gorm.DB, id uint, req tableUserDTO.UpdateTableUserRequest) (models.TableUser, error) {

	if id == 0 {
		return models.TableUser{}, errors.New("invalid table_user_id")
	}

	tableUserData, err := GetTableUser(db, id)
	if err != nil {
		return models.TableUser{}, err
	}

	if err := req.Validate(); err != nil {
		return models.TableUser{}, err
	}

	//Update userTableDTO

	if req.Role != 0 {
		tableUserData.Role = req.Role
	}

	if req.TableID != 0 {
		tableUserData.TableID = req.TableID
	}

	if req.UserID != 0 {
		tableUserData.UserID = req.UserID
	}

	if err := db.Save(&tableUserData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.TableUser{}, errors.New("table_user not found")
		}
		return models.TableUser{}, err
	}
	return tableUserData, nil
}
