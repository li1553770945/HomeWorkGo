package models

import (
	"HomeWorkGo/dao"
	"time"
)

type HomeWorkModel struct {
	ID        int       `json:"id,omitempty" gorm:"primary_key"`
	Name      string    `json:"name"  validate:"required"`
	Desc      string    `json:"desc"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	group     GroupModel
	OwnerID   int
	Owner     UserModel `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
}

func CreateHomework(homework *HomeWorkModel) (err error) {
	err = dao.DB.Create(&homework).Error
	return err
}

func GetHomeworkByID(groupID int) (homework *HomeWorkModel, err error) {
	homework = new(HomeWorkModel)
	err = dao.DB.Where("id = ?", groupID).First(&homework).Error
	if err != nil {
		return nil, err
	}

	return homework, nil
}

func GetHomeWorkByOwnerID(ownerID int, start int, end int) (homework *[]HomeWorkModel, err error) {
	homework = new([]HomeWorkModel)
	err = dao.DB.Where("owner_id = ?", ownerID).Offset(start).Limit(end - start).Find(&homework).Error
	if err != nil {
		return nil, err
	}
	return homework, nil
}

func GetHomeworkNumByOwnerID(ownerID int) (num int64, err error) {
	err = dao.DB.Model(&HomeWorkModel{}).Where("owner_id = ?", ownerID).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, nil
}
func DeleteHomeworkByID(ID int) (err error) {
	err = dao.DB.Where("id=?", ID).Delete(&HomeWorkModel{}).Error
	return
}
