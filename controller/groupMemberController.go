package controller

import (
	"HomeWorkGo/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}

	groupID, _ := strconv.Atoi(json["groupID"].(string))
	joined, err := models.CheckJoin(groupID, uidint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if joined {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您已经加入该组织"})
		return
	}

	err = models.JoinGroup(groupID, uidint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}
