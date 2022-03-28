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
func CreateGroupMember(groupID int, uid int) (err error) {
	group, err := GetGroupByID(groupID)
	if err != nil {
		return err
	}
	err = dao.DB.Model(&group).Association("Members").Append(&UserModel{ID: uid})
	return err
}
func DeleteGroupMember(groupID int, uid int) (err error) {
	err = dao.DB.Where("group_model_id=? AND user_model_id = ?", groupID, uid).Delete(&GroupMemberModel{}).Error
	return err
}

func GetGroupJoined(uid int, start int, end int) (groups *[]GroupModel, err error) {
	if err != nil {
		return nil, err
	}
	groups = new([]GroupModel)
	groupMembers := new([]GroupMemberModel)
	dao.DB.Model(&GroupMemberModel{}).Order("created_at desc").Offset(start).Limit(end-start).Where("user_model_id = ?", uid).Find(&groupMembers)
	err = dao.DB.Model(&groupMembers).Association("Group").Find(&groups)
	for i := 0; i < len(*groups); i++ {
		err = dao.DB.Model(&(*groups)[i]).Select("id,name").Association("Owner").Find(&(*groups)[i].Owner)
		if err != nil {
			return nil, err
		}
	}

	return groups, nil
}

func GetGroupJoinedNum(uid int) (num int64, err error) {
	err = dao.DB.Model(&GroupMemberModel{}).Where("user_model_id = ?", uid).Count(&num).Error
	return num, err
}
