package models

import (
	"HomeWorkGo/dao"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type GroupModel struct {
	ID          int       `json:"id,omitempty" gorm:"primary_key"`
	Name        string    `json:"name,omitempty"  validate:"required"`
	Desc        string    `json:"desc,omitempty"`
	Password    string    `json:"-"  validate:"required"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	AllowCreate bool      `json:"allowCreate"`
	SavePath    string    `json:"-"`
	OwnerID     int
	Owner       UserModel   `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	Members     []UserModel `json:"members,omitempty" gorm:"many2many:group_members;"`
}

func CreateGroup(group *GroupModel) (err error) {
	group.SavePath = filepath.ToSlash(filepath.Join(strconv.Itoa(time.Now().Year()), strconv.Itoa(int(time.Now().Month()))))
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	err = os.MkdirAll(filepath.ToSlash(filepath.Join(dir, group.SavePath)), os.ModePerm)
	if err != nil {
		return err
	}
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
