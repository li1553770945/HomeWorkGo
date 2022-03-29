package models

import (
	"HomeWorkGo/dao"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type HomeWorkModel struct {
	ID                int       `json:"id,omitempty" gorm:"primary_key"`
	Name              string    `json:"name"  validate:"required"`
	Desc              string    `json:"desc"`
	Subject           string    `json:"subject" validate:"required"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	EndTime           time.Time `json:"endtime" validate:"required"`
	CanSubmitAfterEnd bool
	GroupID           int        `validate:"required"`
	Group             GroupModel `gorm:"Foreignkey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	OwnerID           int        `validate:"required"`
	Owner             UserModel  `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	SavePath          string     `json:"-"`
	//Submissions       []SubmissionModel `json:"submissions,omitempty" gorm:"many2many:submissions;"`
	Type string `json:"type"`
}

func CreateHomework(homework *HomeWorkModel) (err error) {
	err = dao.DB.Create(&homework).Error
	if err != nil {
		return err
	}

	homework.SavePath = filepath.ToSlash(filepath.Join(strconv.Itoa(time.Now().Year()), strconv.Itoa(int(time.Now().Month())), strconv.Itoa(homework.ID)))

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(dir)
	err = os.MkdirAll(filepath.ToSlash(filepath.Join(dir, homework.SavePath)), os.ModePerm)
	if err != nil {
		return err
	}
	err = UpdateHomework(homework)
	if err != nil {
		return err
	}

	group, err := GetGroupByID(homework.GroupID)
	if err != nil {
		return err
	}
	for i := 0; i < len(group.Members); i++ {
		err = dao.DB.Create(&SubmissionModel{OwnerID: group.Members[i].ID, HomeworkID: homework.ID, EndTime: homework.EndTime}).Error
		if err != nil {
			return err
		}
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
	if end-start > 100 {
		end = start + 100
	}
	err = dao.DB.Where("owner_id = ?", ownerID).Order("created_at desc").Offset(start).Limit(end - start).Find(&homework).Error
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

func UpdateHomework(homework *HomeWorkModel) (err error) {
	homework_before, err := GetHomeworkByID(homework.ID)
	if homework_before.EndTime != homework.EndTime {
		submissions, err := GetSubmissionsByHomeworkID(homework.ID)
		if err != nil {
			return err
		}
		for i := 0; i < len(*submissions); i++ {
			(*submissions)[i].EndTime = homework.EndTime
			err := UpdateSubmission(&(*submissions)[i])
			if err != nil {
				return err
			}
		}
	}
	err = dao.DB.Save(homework).Error
	return err
}
