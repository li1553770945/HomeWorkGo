package models

import (
	"HomeWorkGo/dao"
	"time"
)

type SubmissionModel struct {
	ID         int       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	HomeworkID int
	Homework   HomeWorkModel `json:"-" gorm:"Foreignkey:HomeworkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	OwnerID    int
	Owner      UserModel `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	FileName   string
	Finish     bool
}

func GetSubmissionsByHomeworkID(homeworkID int) (submissions *[]SubmissionModel, err error) {
	submissions = new([]SubmissionModel)
	err = dao.DB.Where("homework_id = ?", homeworkID).Find(&submissions).Error
	for i := 0; i < len(*submissions); i++ {
		err = dao.DB.Model(&(*submissions)[i]).Select("id,name").Association("Owner").Find(&(*submissions)[i].Owner)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return submissions, nil
}
func GetSubmissionByID(ID int) (submission *SubmissionModel, err error) {
	submission = new(SubmissionModel)

	err = dao.DB.Where("id = ?", ID).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return submission, err
}
