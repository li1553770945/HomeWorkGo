package models

import (
	"HomeWorkGo/dao"
	"time"
)

type GroupModel struct {
	ID          int       `json:"id,omitempty" gorm:"primary_key"`
	Name        string    `json:"name,omitempty"  validate:"required"`
	Desc        string    `json:"desc,omitempty"`
	Password    string    `json:"password"  validate:"required"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	AllowCreate bool      `json:"allowCreate"`
	OwnerID     int
	Owner       UserModel   `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	Members     []UserModel `json:"members,omitempty" gorm:"many2many:group_members;"`
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
	err = dao.DB.Model(&group).Select("id,name").Association("Owner").Find(&group.Owner)
	if err != nil {
		return nil, err
	}
	err = dao.DB.Model(&group).Select("id,name").Association("Members").Find(&group.Members)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func GetGroupsByOwnerID(ownerID int, start int, end int) (groups *[]GroupModel, err error) {
	groups = new([]GroupModel)
	if end-start > 100 {
		end = start + 100
	}
	err = dao.DB.Where("owner_id = ?", ownerID).Order("created_at desc").Offset(start).Limit(end - start).Find(&groups).Error
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

func UpdateGroup(group *GroupModel) (err error) {
	err = dao.DB.Save(group).Error
	return err
}
