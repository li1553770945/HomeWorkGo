package controller

import (
	"HomeWorkGo/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfig(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	uidInt := uid.(int)

	config, err := models.GetConfigByUserID(uidInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": config})
	return
}

func UpdateConfig(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	uidInt := uid.(int)
	config, err := models.GetConfigByUserID(uidInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	jsonData := make(map[string]interface{})
	c.BindJSON(&jsonData)
	if jsonData["config"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	config.Config = jsonData["config"].(string)
	err = models.UpdateConfig(config)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
}
