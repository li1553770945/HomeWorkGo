package controller

import (
	"HomeWorkGo/models"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateHomework(c *gin.Context) {

	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	uidint := uid.(int)

	var homework models.HomeWorkModel
	err := c.BindJSON(&homework)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}

	homework.OwnerID = uidint

	validate := validator.New()
	err = validate.Struct(homework)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	group, err := models.GetGroupByID(homework.GroupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "使用了不存在的小组"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	joined, err := models.CheckJoin(homework.GroupID, uidint)
	if !(group.OwnerID == uidint || (group.AllowCreate && joined)) {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "使用了没有权限使用的小组"})
		return
	}

	err = models.CreateHomework(&homework)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": homework.ID,
		})
		return
	}
}

func GetHomework(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	homeworkID := c.Query("homeworkID")

	if homeworkID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	homeworkIDInt, _ := strconv.Atoi(homeworkID)
	homework, err := models.GetHomeworkByID(homeworkIDInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的作业不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if homework.OwnerID != uid.(int) {
		_, err := models.GetSubmissionByHomeworkAndOwner(uid.(int), homeworkIDInt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{"ode": 4003, "msg": "您没有权限查看"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": homework})
	return
}

func UpdateHomework(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	homeworkID := c.Query("homeworkID")

	if homeworkID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	homeworkIDInt, _ := strconv.Atoi(homeworkID)
	homework, err := models.GetHomeworkByID(homeworkIDInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的作业不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if homework.OwnerID != uid.(int) {
		c.JSON(http.StatusOK, gin.H{"ode": 4003, "msg": "您没有权限执行该操作"})
		return
	}
	err = c.BindJSON(&homework)
	if err != nil || homework.ID != homeworkIDInt {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	err = models.UpdateHomework(homework)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func DeleteHomework(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	jsonData := make(map[string]interface{})
	err := c.BindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	if jsonData["homeworkID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	homeworkID := int(jsonData["homeworkID"].(float64))
	homework, err := models.GetHomeworkByID(homeworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的作业不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if homework.OwnerID != uid.(int) {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限执行该操作"})
		return
	}
	err = models.DeleteHomeworkByID(homeworkID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func GetHomeworkCreated(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	start, end := c.Query("start"), c.Query("end")

	if start == "" || end == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	ownerIDint := uid.(int)
	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)

	homework, err := models.GetHomeworkByOwnerID(ownerIDint, startInt, endInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": homework})
	return
}

func GetHomeworkCreatedNum(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	ownerIDint := uid.(int)

	num, err := models.GetHomeworkNumByOwnerID(ownerIDint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
	return
}
