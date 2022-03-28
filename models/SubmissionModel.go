package models

import (
	"HomeWorkGo/dao"
	"github.com/fatih/structs"
	"time"
)

type SubmissionModel struct {
	ID         int       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	EndTime    time.Time `json:"-"`
	SubmitAt   time.Time
	HomeworkID int
	Homework   HomeWorkModel `json:"-" gorm:"Foreignkey:HomeworkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	OwnerID    int
	Owner      UserModel `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`

	FileName string
	Finished bool
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

func GetHomeworkJoinedByOwnerId(ownerID int, start int, end int) (results *[]map[string]interface{}, err error) {
	if end-start > 100 {
		end = start + 100
	}
	results = new([]map[string]interface{})
	submissions := new([]SubmissionModel)
	err = dao.DB.Order("created_at desc").Where("owner_id = ?", ownerID).Offset(start).Limit(end - start).Find(submissions).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(*submissions); i++ {
		homework, err := GetHomeworkByID((*submissions)[i].HomeworkID)
		if err != nil {
			return nil, err
		}

		m := structs.Map(homework)
		delete(m, "Group")
		delete(m, "GroupID")
		delete(m, "OwnerID")
		delete(m, "Submissions")
		delete(m, "Owner")
		owner := structs.Map(homework.Owner)
		delete(owner, "CreatedAt")
		delete(owner, "LastLogin")
		delete(owner, "Password")
		delete(owner, "Status")
		delete(owner, "Username")
		delete(owner, "Validation")
		m["Owner"] = owner
		m["Finished"] = (*submissions)[i].Finished
		*results = append(*results, m)
	}
	return results, nil

}

func GetHomeworkNotFinishedByOwnerId(ownerID int, start int, end int) (results *[]map[string]interface{}, err error) {
	if end-start > 100 {
		end = start + 100
	}
	results = new([]map[string]interface{})
	submissions := new([]SubmissionModel)
	err = dao.DB.Order("end_time desc").Where("owner_id = ? AND finished = ? ", ownerID, 0).Offset(start).Limit(end - start).Find(submissions).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(*submissions); i++ {
		homework, err := GetHomeworkByID((*submissions)[i].HomeworkID)
		if err != nil {
			return nil, err
		}

		m := structs.Map(homework)
		delete(m, "Group")
		delete(m, "GroupID")
		delete(m, "OwnerID")
		delete(m, "Submissions")
		delete(m, "Owner")
		owner := structs.Map(homework.Owner)
		delete(owner, "CreatedAt")
		delete(owner, "LastLogin")
		delete(owner, "Password")
		delete(owner, "Status")
		delete(owner, "Username")
		delete(owner, "Validation")
		m["Owner"] = owner
		m["Finished"] = (*submissions)[i].Finished
		*results = append(*results, m)
	}
	return results, nil

}

func GetSubmissionByHomeworkAndOwner(ownerID int, homewrokID int) (submission *SubmissionModel, err error) {
	submission = new(SubmissionModel)
	err = dao.DB.Model(&SubmissionModel{}).Where("owner_id = ? AND homework_id = ?", ownerID, homewrokID).First(submission).Error
	return submission, err
}
func GetHomeworkJoinedNumByOwnerId(ownerID int) (num int64, err error) {

	err = dao.DB.Model(&SubmissionModel{}).Where("owner_id = ?", ownerID).Count(&num).Error
	if err != nil {
		return 0, err
	}
	return num, err
}

func UpdateSubmission(submission *SubmissionModel) (err error) {
	err = dao.DB.Save(submission).Error
	return err
}
