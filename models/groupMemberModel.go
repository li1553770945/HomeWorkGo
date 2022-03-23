package models

import (
	"HomeWorkGo/dao"
	"time"
)

type GroupMemberModel struct {
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UserModelID  int        `gorm:"primaryKey;"`
	GroupModelID int        `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User         UserModel  `gorm:"Foreignkey:UserModelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group        GroupModel `gorm:"Foreignkey:GroupModelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func CheckJoin(groupID int, uid int) (joined bool, err error) {
	var num int64
	err = dao.DB.Model(&GroupMemberModel{}).Where("group_model_id = ? AND user_model_id = ?", groupID, uid).Count(&num).Error
	return num != 0, err
}
func JoinGroup(groupID int, uid int) (err error) {
	group, err := GetGroupByID(groupID)
	if err != nil {
		return err
	}
	err = dao.DB.Model(&group).Association("Members").Append(&UserModel{ID: uid})
	return err
}