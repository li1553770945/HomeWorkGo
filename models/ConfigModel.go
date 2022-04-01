package models

import (
	"HomeWorkGo/dao"
)

type ConfigModel struct {
	ID      int    `json:"id,omitempty" gorm:"primary_key"`
	OwnerID int    `gorm:"index"`
	Config  string `json:"config"`
}

func CreateConfig(ownerId int) (err error) {
	err = dao.DB.Create(&ConfigModel{OwnerID: ownerId}).Error
	return err
}
func UpdateConfig(config *ConfigModel) (err error) {
	err = dao.DB.Save(&config).Error
	return err
}

func GetConfigByUserID(userID int) (config *ConfigModel, err error) {
	config = new(ConfigModel)
	err = dao.DB.Model(&ConfigModel{}).Where("owner_id = ?", userID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return config, nil
}
