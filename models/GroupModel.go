package models

import (
	"HomeWorkGo/dao"
	"time"
)

type GroupModel struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	OwnerID   int
	Owner     UserModel   `gorm:"Foreignkey:OwnerID"`
	Members   []UserModel `gorm:"many2many:group_members;"`
}

func CreateGroup(group *GroupModel) (err error) {
	err = dao.DB.Create(&group).Error
	return err
}

func GetGroupByID(groupID int) (group *GroupModel, err error) {
	group = new(GroupModel)
	err = dao.DB.Where("id = ?", groupID).First(&group).Error
	if err != nil {
		return nil, err
	}
	err = dao.DB.Model(&group).Association("owner").Error
	if err != nil {
		return nil, err
	}
	return group, nil
}

func GetGroupsByOwnerID(ownerID int, start int, end int) (groups *[]GroupModel, err error) {
	groups = new([]GroupModel)
	err = dao.DB.Where("owner_id = ?", ownerID).Offset(start).Limit(end - start).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func GetGroupNumByOwnerID(ownerID int) (num int64, err error) {
	err = dao.DB.Model(&GroupModel{}).Where("owner_id = ?", ownerID).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}
func DeleteGroupByID(groupID int) (err error) {
	err = dao.DB.Where("id=?", groupID).Delete(&GroupModel{}).Error
	return
}
