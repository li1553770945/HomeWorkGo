package models

//
//import (
//	"HomeWorkGo/dao"
//	"time"
//)
//
//type HomeWorkModel struct {
//	ID        int       `json:"id,omitempty" gorm:"primary_key"`
//	Name      string    `json:"name"  validate:"required"`
//	Desc      string    `json:"desc"`
//	Subject   string    `json:"subject"`
//	CreatedAt time.Time `gorm:"autoCreateTime"`
//	group     GroupModel
//	OwnerID   int
//	Owner     UserModel `json:"owner,omitempty" gorm:"Foreignkey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"  validate:"-"`
//}
//
//func CreateHomework(homework *HomeWorkModel) (err error) {
//	err = dao.DB.Create(&homework).Error
//	return err
//}
//
//func GetHomeworkByID(groupID int) (group *GroupModel, err error) {
//	group = new(GroupModel)
//	err = dao.DB.Where("id = ?", groupID).First(&group).Error
//	if err != nil {
//		return nil, err
//	}
//	err = dao.DB.Model(&group).Select("id,name").Association("Owner").Find(&group.Owner)
//	err = dao.DB.Model(&group).Select("id,name").Association("Members").Find(&group.Members)
//	if err != nil {
//		return nil, err
//	}
//	return group, nil
//}
//
//func GetGroupsByOwnerID(ownerID int, start int, end int) (groups *[]GroupModel, err error) {
//	groups = new([]GroupModel)
//	err = dao.DB.Where("owner_id = ?", ownerID).Offset(start).Limit(end - start).Find(&groups).Error
//	if err != nil {
//		return nil, err
//	}
//	return groups, nil
//}
//
//func GetGroupNumByOwnerID(ownerID int) (num int64, err error) {
//	err = dao.DB.Model(&GroupModel{}).Where("owner_id = ?", ownerID).Count(&num).Error
//	if err != nil {
//		return 0, err
//	}
//	return num, nil
//}
//func DeleteGroupByID(groupID int) (err error) {
//	err = dao.DB.Where("id=?", groupID).Delete(&GroupModel{}).Error
//	return
//}
