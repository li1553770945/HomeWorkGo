package controller

import (
	"HomeWorkGo/models"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func JoinGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	uidint := uid.(int)
	json := make(map[string]interface{})
	c.BindJSON(&json)
	if json["groupID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	groupID := int(json["groupID"].(float64))
	joined, err := models.CheckJoin(groupID, uidint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if joined {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您已经加入该组织"})
		return
	}
	group, err := models.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的小组不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if group.Password != json["password"] {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "加入密码错误"})
		return
	}

	err = models.CreateGroupMember(groupID, uidint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func QuitGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	uidint := uid.(int)
	json := make(map[string]interface{})
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	if json["groupID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	groupID := int(json["groupID"].(float64))
	joined, err := models.CheckJoin(groupID, uidint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if !joined {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有加入该组织"})
		return
	}

	err = models.DeleteGroupMember(groupID, uidint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func GetGroupsJoined(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	ownerIDint := uid.(int)
	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)

	groups, err := models.GetGroupJoined(ownerIDint, startInt, endInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": groups})
	return
}

func GetMyGroupsJoinedNum(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	ownerIDint := uid.(int)

	num, err := models.GetGroupJoinedNum(ownerIDint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
	return
}
