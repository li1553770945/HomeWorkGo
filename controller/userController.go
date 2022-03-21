package controller

import (
	"HomeWorkGo/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func Register(c *gin.Context) {

	var user models.UserModel
	user.Status = true
	validate := validator.New()

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": err.Error()})
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
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": err.Error()})
		return
	}
	err = validate.Struct(loginRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": err.Error()})
		return
	}

	var user *models.UserModel
	fmt.Println(loginRequest.Username)
	user, err = models.GetUserByUserName(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "用户名或密码错误"})
		return
	}
	if user.Password != loginRequest.Password {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "用户名或密码错误"})
		return
	}

	session := sessions.Default(c)
	session.Set("uid", user.ID)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"code": 0})

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
