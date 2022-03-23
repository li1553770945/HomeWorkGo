package controller

import (
	"HomeWorkGo/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
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
	c.BindJSON(&json)
	if json["groupID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
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

	groups, err := models.GetGroupJoined(ownerIDint, startint, endint)
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
