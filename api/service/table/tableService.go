package table

import (
	"errors"
	"github.com/GarotoCowboy/vttProject/api/dto/tableDTO"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/GarotoCowboy/vttProject/api/utils"
	"gorm.io/gorm"
)

func CreateTable(db *gorm.DB, req tableDTO.CreateTableRequest) (models.Table, error) {

	//Validate if request bod is valid
	if err := req.Validate(); err != nil {
		return models.Table{}, err
	}

	generatedInviteLink, err := createAndVerifyLinkCodeNotExists(db)
	if err != nil {
		return models.Table{}, err
	}

	var table = models.Table{
		Name:     req.Name,
		Password: req.Password,
		OwnerID:  req.OwnerID,
	}
	table.InviteLink = generatedInviteLink

	//hash password to bcrypt
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.Table{}, err
	}

	table.Password = string(hashedPassword)

	//Save Table
	if err := db.Create(&table).Error; err != nil {
		return models.Table{}, err
	}

	ownerMember := models.TableUser{
		TableID: table.ID,
		UserID:  req.OwnerID,
		Role:    consts.Role(2),
	}

	if err := db.Create(&ownerMember).Error; err != nil {
		return models.Table{}, err
	}

	var fullTable models.Table
	if err := db.Preload("Owner").First(&fullTable, table.ID).Error; err != nil {
		return models.Table{}, err
	}
	return fullTable, nil
}

func DeleteTable(db *gorm.DB, id uint) (models.Table, error) {

	//Verify if id is valid
	if id == 0 {
		return models.Table{}, errors.New("invalid table_id")
	}

	var table = models.Table{}

	//Preload all members registred at table
	if err := db.Preload("Members").First(&table, id).Error; err != nil {
		return models.Table{}, err
	}
	//Select all members
	if err := db.Select("Members").Delete(&table).Error; err != nil {
		return models.Table{}, err
	}
	//Delete table and members with inputted id
	if err := db.Delete(&table).Error; err != nil {
		return models.Table{}, err
	}
	return table, nil
}

func GetTable(db *gorm.DB, id uint) (models.Table, error) {
	if id <= 0 {
		return models.Table{}, errors.New("invalid table_id")
	}

	var table = models.Table{}

	if err := db.Preload("Owner").Where("id = ?", id).First(&table).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Table{}, errors.New("table not found")
		}
		return models.Table{}, err
	}

	return table, nil
}

func ListTables(db *gorm.DB) ([]models.Table, error) {

	var tables []models.Table

	if err := db.Preload("Owner").Find(&tables).Error; err != nil {
		return tables, err
	}
	return tables, nil
}

func UpdateTable(db *gorm.DB, id uint, req tableDTO.UpdateTableRequest) (models.Table, error) {

	if id == 0 {
		return models.Table{}, errors.New("invalid table_id")
	}

	tableData, err := GetTable(db, id)
	if err != nil {
		return models.Table{}, err
	}

	if err := req.Validate(); err != nil {
		return models.Table{}, err
	}

	//Update userDTO

	if req.Name != "" {
		tableData.Name = req.Name
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return models.Table{}, err
		}

		tableData.Password = string(hashedPassword)
	}

	if err := db.Save(&tableData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Table{}, errors.New("table not found")
		}
		return models.Table{}, err
	}
	return tableData, nil
}

func createAndVerifyLinkCodeNotExists(db *gorm.DB) (string, error) {
	var count int64
	var link string
	for {
		link = utils.StringWithCharset(10)

		err := db.
			Model(&models.Table{}).
			Where("invite_link = ?", link).
			Count(&count).Error
		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}
	}
	return link, nil

}
