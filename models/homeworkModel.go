package models

import (
	"HomeWorkGo/dao"
	"fmt"
	"time"
)

type HomeWorkModel struct {
	ID                int               `json:"id,omitempty" gorm:"primary_key"`
	Name              string            `json:"name"  validate:"required"`
	Desc              string            `json:"desc"`
	Subject           string            `json:"subject" validate:"required"`
	CreatedAt         time.Time         `gorm:"autoCreateTime"`
	EndTime           time.Time         `json:"endTime" validate:"required"`
	CanSubmitAfterEnd bool              `validate:"required"`
	GroupID           int               `validate:"required"`
	Group             GroupModel        `gorm:"Foreignkey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	OwnerID           int               `validate:"required"`
	Owner             UserModel         `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	Submissions       []SubmissionModel `json:"submissions,omitempty" gorm:"many2many:submissions;"`
}

func CreateHomework(homework *HomeWorkModel) (err error) {
	err = dao.DB.Create(&homework).Error
	group, err := GetGroupByID(homework.GroupID)
	if err != nil {
		return err
	}

	for i := 0; i < len(group.Members); i++ {
		fmt.Println(group.Members[i].ID, homework.ID)
		err = dao.DB.Model(&homework).Association("Submissions").Append(&SubmissionModel{OwnerID: group.Members[i].ID, HomeworkID: homework.ID})
	}
	return err
}

func GetHomeworkByID(groupID int) (homework *HomeWorkModel, err error) {
	homework = new(HomeWorkModel)
	err = dao.DB.Where("id = ?", groupID).First(&homework).Error
	if err != nil {
		return nil, err
	}
	err = dao.DB.Model(&homework).Select("id,name").Association("Owner").Find(&homework.Owner)
	if err != nil {
		return nil, err
	}
	return homework, nil
}

func GetHomeworkByOwnerID(ownerID int, start int, end int) (homework *[]HomeWorkModel, err error) {
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
