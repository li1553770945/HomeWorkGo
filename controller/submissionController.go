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

func GetSubmissionsByHomeworkId(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	uidInt := uid.(int)
	homeworkID := c.Query("homeworkID")

	if homeworkID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}

	homeworkIDInt, _ := strconv.Atoi(homeworkID)
	homework, err := models.GetHomeworkByID(homeworkIDInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的作业不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if homework.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您不是该作业的创建者，无权查看完成情况"})
		return
	}

	submissions, err := models.GetSubmissionsByHomeworkID(homeworkIDInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": submissions})
	return
}
func GetSubmissionFileById(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	uidInt := uid.(int)
	submissionID := c.Query("submissionID")

	if submissionID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 2001, "msg": "请求参数错误"})
		return
	}

	submissionIDInt, _ := strconv.Atoi(submissionID)
	submission, err := models.GetSubmissionByID(submissionIDInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的作业不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

	if submission.OwnerID != uidInt {
		homework, _ := models.GetHomeworkByID(submission.HomeworkID)
		if homework.OwnerID != uid {
			c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限查看该内容"})
			return
		}
	}
	if !submission.Finish {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "还没有完成该作业"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})
	return
}
