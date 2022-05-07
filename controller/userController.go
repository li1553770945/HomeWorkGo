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
	"time"
)

func Register(c *gin.Context) {

	var user models.UserModel
	user.Status = 1
	validate := validator.New()
	data, err := c.GetRawData()
	if err != nil {
		fmt.Println(err.Error())
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	err = c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	jsonData := make(map[string]interface{})
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if jsonData["password"] == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "提交参数错误"})
		return
	}

	models.SetUserPasswd(&user, jsonData["password"].(string))

	_, err = models.GetUserByUserName(user.Username)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "用户名已经存在"})
		return
	}

	err = models.CreateUser(&user)
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

func Login(c *gin.Context) {

	var loginRequest models.LoginRequest

	validate := validator.New()

	err := c.BindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	err = validate.Struct(loginRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}

	var user *models.UserModel
	user, err = models.GetUserByUserName(loginRequest.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "用户名或密码错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if !models.CheckUserPasswd(user, loginRequest.Password) {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "用户名或密码错误"})
		return
	}
	if user.Status == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "抱歉，您的账户目前不可用"})
	}
	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()
	user.LastLogin = time.Now()
	models.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})

}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("uid")
	session.Save()
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func Me(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	user, err := models.GetUserById(uid.(int))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": user})
	return
}

func NewValidation(c *gin.Context) {
	//email := c.Query("email")
	//user =
	//token := utils.GetToken()
	//dao.RDB.Set("resetPassword"+email, token, 300*time.Second)
	//c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func UpdateUser(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	user, err := models.GetUserById(uid.(int))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	jsonData := make(map[string]interface{})
	err = c.BindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	name, exist := jsonData["name"]
	if exist {
		user.Name = name.(string)
	}

	username, exist := jsonData["username"]
	if exist {
		user.Username = username.(string)
	}

	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}
	err = models.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}

func ResetPassword(c *gin.Context) {
	jsonData := make(map[string]interface{})
	err := c.BindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	username := jsonData["username"]
	validation := jsonData["validation"]
	password := jsonData["password"]
	if username == nil || validation == nil || password == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "参数不合法"})
		return
	}
	rUser, err := models.GetUserByUserName(username.(string))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "用户不存在或验证码错误"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if rUser.ResetCnt == 3 {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "用户不存在或验证码错误"})
		return
	}
	if len(rUser.Validation) != 8 || rUser.Validation != validation.(string) {
		rUser.ResetCnt += 1
		models.UpdateUser(rUser)
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "用户不存在或验证码错误"})
		return

	}
	err = models.ResetPassword(rUser, password.(string))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return

}
