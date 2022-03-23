package models

import (
	"HomeWorkGo/dao"
	"time"
)

type GroupMemberModel struct {
	ID        int       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	OwnerID   int       `gorm:"index"`
	GroupID   int       `gorm:"index"`
}

func CreateGroupMember(groupMember *GroupMemberModel) (err error) {
	err = dao.DB.Create(&groupMember).Error
	return err
}

func DeleteGroupMember(groupMember *GroupMemberModel) (err error) {
	err = dao.DB.Create(&groupMember).Error
	return err
}

func GetMyGroups(groupMember *GroupMemberModel) (err error) {
	err = dao.DB.Create(&groupMember).Error
	return err
}
