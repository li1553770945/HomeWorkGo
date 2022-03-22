package models

import (
	"HomeWorkGo/dao"
	"crypto/md5"
	"encoding/hex"
	"time"
)

type UserModel struct {
	ID         int       `json:"id" gorm:"primary_key"`
	Username   string    `gorm:"unique_index"`
	Name       string    `json:"name"`
	Password   string    `json:"-"`
	Validation string    `json:"-"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	LastLogin  time.Time
}

func CreateUser(user *UserModel) (err error) {
	err = dao.DB.Create(&user).Error
	return err
}

func GetUserByUserName(username string) (user *UserModel, err error) {
	user = new(UserModel)
	if err = dao.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	return
}

func GetUserById(id int) (user *UserModel, err error) {
	user = new(UserModel)
	if err = dao.DB.Debug().Where("id=?", id).First(user).Error; err != nil {
		return nil, err
	}
	return
}
func Encrypt(passwd string) string {
	m := md5.New()
	m.Write([]byte(passwd))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
func SetUserPasswd(user *UserModel, passwd string) {
	user.Password = Encrypt(passwd)
}
func CheckUserPasswd(user *UserModel, passwd string) bool {
	return Encrypt(passwd) == user.Password
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
