package controller

import (
	"HomeWorkGo/dao"
	"HomeWorkGo/models"
	"HomeWorkGo/utils"
	"archive/zip"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
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
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
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
	homeworkID := c.Query("homeworkID")
	if submissionID == "" && homeworkID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	homeworkIDInt, _ := strconv.Atoi(homeworkID)
	submission := &models.SubmissionModel{}
	if submissionID == "" {
		submissionTemp, err := models.GetSubmissionByHomeworkAndOwner(uidInt, homeworkIDInt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限查看该内容"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
			return
		}
		submission = submissionTemp
	} else {
		submissionIDInt, _ := strconv.Atoi(submissionID)
		submissionTemp, err := models.GetSubmissionByID(submissionIDInt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "请求的提交不存在"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
			return
		}
		submission = submissionTemp

	}

	if submission.OwnerID != uidInt {
		homework, _ := models.GetHomeworkByID(submission.HomeworkID)
		if homework.OwnerID != uid {
			c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您没有权限查看该内容"})
			return
		}
	}
	if !submission.Finished {
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "还没有完成该作业"})
		return
	}

	homework, _ := models.GetHomeworkByID(submission.HomeworkID)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fullPath := filepath.ToSlash(filepath.Join(dir, homework.SavePath, submission.FileName))
	submissionIDString := strconv.Itoa(submission.ID)
	token := "file" + submissionIDString + utils.GetToken()
	err := dao.RDB.Set(token, fullPath, 60*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"token": token}})
	return
}

func GetHomeworkJoinedNum(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	ownerIDint := uid.(int)

	num, err := models.GetHomeworkJoinedNumByOwnerId(ownerIDint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
	return
}

func GetHomeworkJoined(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	start, end := c.Query("start"), c.Query("end")

	if start == "" || end == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	ownerIDint := uid.(int)
	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)

	homework, err := models.GetHomeworkJoinedByOwnerId(ownerIDint, startInt, endInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": homework})
	return
}

func GetHomeworkNotFinished(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	start, end := c.Query("start"), c.Query("end")

	if start == "" || end == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	ownerIDint := uid.(int)
	startInt, _ := strconv.Atoi(start)
	endInt, _ := strconv.Atoi(end)

	homework, err := models.GetHomeworkNotFinishedByOwnerId(ownerIDint, startInt, endInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": homework})
	return
}

func GetHomeworkNotFinishedNum(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	ownerIDint := uid.(int)

	num, err := models.GetHomeworkNotFinishedNumByOwnerId(ownerIDint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
	return
}

func Submit(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}
	uidInt := uid.(int)
	homeworkID := c.Query("homeworkID")

	if homeworkID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}

	homeworkIDInt, _ := strconv.Atoi(homeworkID)
	submission, err := models.GetSubmissionByHomeworkAndOwner(uidInt, homeworkIDInt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "您没有参与该次作业"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	submission.SubmitAt = time.Now()
	user, err := models.GetUserById(uidInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	homework, err := models.GetHomeworkByID(homeworkIDInt)
	if time.Now().After(homework.EndTime) && !homework.CanSubmitAfterEnd {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "抱歉，该作业提交截至时间已过"})
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	submission.FileName = user.Username + user.Name + path.Ext(file.Filename)
	full_path := filepath.ToSlash(filepath.Join(dir, homework.SavePath, submission.FileName))
	err = c.SaveUploadedFile(file, full_path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	submission.Finished = true
	err = models.UpdateSubmission(submission)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0})

}

func ExportThread(homeworkID string, name string, savePath string) {
	token := "export" + homeworkID + utils.GetToken()
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	zipPath := filepath.ToSlash(filepath.Join(dir, "export", name+time.Now().Format("20060102150405")+".zip"))

	f, err := ioutil.ReadDir(filepath.ToSlash(filepath.Join(dir, savePath)))
	if err != nil {
		fmt.Println("压缩目录读取错误")
		dao.RDB.Set("error"+homeworkID, "压缩目录读取错误", 3*time.Minute)
		dao.RDB.Del("export" + homeworkID)
		return
	}

	fZip, _ := os.Create(zipPath)
	w := zip.NewWriter(fZip)

	for _, file := range f {
		fw, err := w.Create(file.Name())
		if err != nil {
			fmt.Println("创建失败", err)
			dao.RDB.Set("error"+homeworkID, "添加文件错误", 3*time.Minute)
			dao.RDB.Del("export" + homeworkID)
			return
		}

		fileContent, err := ioutil.ReadFile(filepath.ToSlash(filepath.Join(dir, savePath, file.Name())))
		if err != nil {
			fmt.Println("读取文件错误")
			dao.RDB.Set("error"+homeworkID, "读取文件错误", 3*time.Minute)
			dao.RDB.Del("export" + homeworkID)
			return
		}
		_, err = fw.Write(fileContent)
		if err != nil {
			fmt.Println("写文件错误")
			dao.RDB.Set("error"+homeworkID, "写入压缩文件错误", 3*time.Minute)
			dao.RDB.Del("export" + homeworkID)
			return
		}
	}
	err = w.Close()
	if err != nil {
		return
	}
	fmt.Printf("打包完成")
	dao.RDB.Set("export"+homeworkID, token, 3*time.Minute)
	dao.RDB.Set(token, zipPath, 3*time.Minute)
	time.Sleep(30 * time.Minute)
	err = os.Remove(zipPath)
	if err != nil {
		return
	}
}

func Export(c *gin.Context) {
	session := sessions.Default(c)
	uid := session.Get("uid")

	if uid == nil {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您还未登录，请先登录"})
		return
	}

	homeworkID := c.Query("homeworkID")

	if homeworkID == "" {
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
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
	uidInt := uid.(int)
	if homework.OwnerID != uidInt {
		c.JSON(http.StatusOK, gin.H{"code": 4003, "msg": "您不是该作业的创建者，无权进行导出操作"})
		return
	}

	erro, err := dao.RDB.Get("error" + homeworkID).Result()
	if err == nil {
		dao.RDB.Del("error" + homeworkID)
		c.JSON(http.StatusOK, gin.H{"code": 5001, "data": erro})
		return
	}
	token, err := dao.RDB.Get("export" + homeworkID).Result()
	if err == nil { //没有错误，两种情况：导出中、导出完成
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"done": false, "token": ""}})
			return
		} else {

			c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"done": true, "token": token}})
			return
		}

	} else {
		if err == redis.Nil {

			dao.RDB.Set("export"+homeworkID, "", 3*time.Hour)
			go ExportThread(homeworkID, homework.Name, homework.SavePath)
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"done": false, "token": ""}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}

}

func Download(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.Status(404)
		c.JSON(http.StatusOK, gin.H{"code": 4001, "msg": "请求参数错误"})
		return
	}
	fullPath, err := dao.RDB.Get(token).Result()
	if err != nil {
		if err == redis.Nil {
			c.Status(404)
			c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "token不存在或已经过期"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 5001, "msg": err.Error()})
		return
	}
	if fullPath == "" {
		c.Status(404)
		c.JSON(http.StatusOK, gin.H{"code": 4004, "msg": "token不存在或已经过期"})
		return
	}

	_, fileName := filepath.Split(fullPath)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	c.File(fullPath)

	return
}
