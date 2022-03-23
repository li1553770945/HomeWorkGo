package controller

import (
	"HomeWorkGo/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func CreateGroup(c *gin.Context) {

	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	uidint := uid.(int)

	var group models.GroupModel

	validate := validator.New()
	err := validate.Struct(group)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}

	err = c.BindJSON(&group)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	group.OwnerID = uidint
	err = models.CreateGroup(&group)

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

func GetGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	json := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&json)
	if json["groupID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}

	groupID := int(json["groupID"].(float64))
	group, err := models.GetGroupByID(groupID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": group})
	return
}

func DeleteGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	json := make(map[string]interface{})
	c.BindJSON(&json)
	if json["groupID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}

	groupID := int(json["groupID"].(float64))
	group, err := models.GetGroupByID(groupID)
	if err != nil {
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

func GetGroupsCreated(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	json := make(map[string]interface{}) //注意该结构接受的内容
	c.BindJSON(&json)
	start := json["start"]
	end := json["end"]
	if start == nil || end == nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}
	ownerIDint := uid.(int)
	startint := int(start.(float64))
	endint := int(end.(float64))

	groups, err := models.GetGroupsByOwnerID(ownerIDint, startint, endint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": groups})
	return
}

func GetGroupsCreatedNum(c *gin.Context) {
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
