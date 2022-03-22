package models

import (
	"HomeWorkGo/dao"
	"fmt"
)

type GroupModel struct {
	ID      int    `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	OwnerID int
	Owner   UserModel `gorm:"Foreignkey:OwnerID"`
}

func CreateGroup(group *GroupModel) (err error) {
	err = dao.DB.Create(&group).Error
	return err
}

func GetGroupByID(groupID int) (group *GroupModel, err error) {
	group = new(GroupModel)
	fmt.Println(groupID)
	err = dao.DB.Where("id = ?", groupID).First(&group).Error
	if err != nil {
		return nil, err
	}
	dao.DB.Model(&group).Association("owner").Find(&group.Owner)
	return group, nil
}
