package models

import (
	"HomeWorkGo/dao"
	"github.com/go-playground/validator/v10"
	"time"
)

type UserModel struct {
	ID         int    `json:"id"`
	Username   string `json:"username" validate:"required"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	Validation string `json:"validation"`
	Status     bool   `json:"status"`
}

func timing(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

func CreateUser(user *UserModel) (err error) {
	err = dao.DB.Create(&user).Error
	return err
}

func GetUserByUserName(username string) (user *UserModel, err error) {
	if err = dao.DB.Find(&user).Error; err != nil {
		return nil, err
	}
	return
}

func GetUserById(id string) (user *UserModel, err error) {
	user = new(UserModel)
	if err = dao.DB.Debug().Where("id=?", id).First(user).Error; err != nil {
		return nil, err
	}
	return
}

func UpdateUser(user *UserModel) (err error) {
	err = dao.DB.Save(user).Error
	return err
}

func DeleteUser(id string) (err error) {
	err = dao.DB.Where("id=?", id).Delete(&Todo{}).Error
	return
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
