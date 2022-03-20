package controller

import (
	"HomeWorkGo/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func Register(c *gin.Context) {

	var user models.UserModel
	validate := validator.New()

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "error": err.Error()})
		return
	}

	err = validate.Struct(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "error": err.Error()})
		return
	}

	err = models.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
		})
		return
	}
}

func Login(c *gin.Context) {

	var loginRequest models.LoginRequest

	validate := validator.New()

	err := c.BindJSON(&loginRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "error": err.Error()})
		return
	}
	err = validate.Struct(loginRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "error": err.Error()})
		return
	}

	var user *models.UserModel
	user, err = models.GetUserByUserName(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "error": err.Error()})
		return
	}
	if user.Password == loginRequest.Password {
		c.JSON(http.StatusOK, gin.H{"code": 0})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "error": "用户名或密码错误"})
		return
	}

}

func Logout(c *gin.Context) {

}

func Me(c *gin.Context) {

}
