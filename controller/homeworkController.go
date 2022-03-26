package controller

import (
	"HomeWorkGo/models"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
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

	err = models.CreateHomework(&homework)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
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
	jsonData := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&jsonData)
	if err != nil {
		return
	}
	if jsonData["ID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}

	homeworkID := int(jsonData["ID"].(float64))
	homework, err := models.GetGroupByID(homeworkID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的作业不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": homework})
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
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}
	if jsonData["groupID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}

	groupID := int(jsonData["groupID"].(float64))
	group, err := models.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的小组不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if group.OwnerID != uid.(int) {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限执行该操作"})
		return
	}
	err = models.DeleteGroupByID(groupID)
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
	jsonData := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误" + err.Error()})
		return
	}
	start := jsonData["start"]
	end := jsonData["end"]
	if start == nil || end == nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}
	ownerIDint := uid.(int)
	startInt := int(start.(float64))
	endInt := int(end.(float64))

	groups, err := models.GetGroupsByOwnerID(ownerIDint, startInt, endInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": groups})
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

	num, err := models.GetGroupNumByOwnerID(ownerIDint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
	return
}
