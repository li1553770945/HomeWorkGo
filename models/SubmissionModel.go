package models

import "time"

type SubmissionModel struct {
	ID        int       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	OwnerID   int
	Owner     UserModel `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
	FileName  string
}
