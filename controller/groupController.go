package controller

import (
	"HomeWorkGo/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
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
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	err = c.BindJSON(&group)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	group.OwnerID = uidint
	jsondata := make(map[string]interface{})
	err = json.Unmarshal(data, &jsondata)
	group.Password = jsondata["password"].(string)
	validate := validator.New()
	err = validate.Struct(group)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}

	err = models.CreateGroup(&group)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": group.ID,
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

	groupID := c.Query("groupID")
	if groupID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	groupIDInt, _ := strconv.Atoi(groupID)
	group, err := models.GetGroupByID(groupIDInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的小组不存在"})
			return
		}
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
	jsonData := make(map[string]interface{})
	err := c.BindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	if jsonData["groupID"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
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

func GetGroupsCreated(c *gin.Context) {
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
	fmt.Println(startInt, endInt)
	groups, err := models.GetGroupsByOwnerID(ownerIDint, startInt, endInt)
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

func UpdateGroup(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	groupID := c.Query("groupID")

	if groupID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	groupIDInt, _ := strconv.Atoi(groupID)
	group, err := models.GetGroupByID(groupIDInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的小组不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if group.OwnerID != uid.(int) {
		c.JSON(http.StatusOK, gin.H{"ode": 4003, "msg": "您没有权限执行该操作"})
		return
	}

	err = c.BindJSON(&group)
	if err != nil || group.ID != groupIDInt {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	err = models.UpdateGroup(group)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}
