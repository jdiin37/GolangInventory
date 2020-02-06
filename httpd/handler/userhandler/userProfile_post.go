package userhandler

import (
	"database/sql"
	"inventory/httpd/middleware"
	"inventory/httpd/model"
	"inventory/kit"
	"inventory/platform/user"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type userProfilePostResquest struct {
	Id       int            `form:"id"`
	Email    string         `form:"email" binding:"email"`
	Name     string         `form:"name"`
	Password string         `form:"password" `
	Avatar   sql.NullString `form:"avatar"`
}

func UserProfilePost(userUpdater user.Updater) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := userProfilePostResquest{}

		result := &model.Result{
			Code:    200,
			Message: "更新成功",
			Data:    nil,
		}

		if err := c.ShouldBind(&requestBody); err != nil {
			c.HTML(http.StatusOK, "error.gohtml", gin.H{
				"error": err.Error(),
			})
			//log.Panicln("綁定發生錯誤 ", err.Error())
			result.Code = http.StatusBadRequest
			result.Message = "資料格式有誤"
		}
		file, e := c.FormFile("avatar-file")
		if e != nil {
			// c.HTML(http.StatusOK, "error.gohtml", gin.H{
			// 	"error": e,
			// })
			//log.Panicln("文件上傳錯誤", e.Error())
		} else {
			path := kit.RootPath()
			path = filepath.Join(path, "avatar")
			e = os.MkdirAll(path, os.ModePerm)
			if e != nil {
				c.HTML(http.StatusOK, "error.gohtml", gin.H{
					"error": e,
				})
				log.Panicln("無法創建文件夾", e.Error())
			}
			fileName := strconv.FormatInt(time.Now().Unix(), 10) + file.Filename
			e = c.SaveUploadedFile(file, path+"/"+fileName)
			if e != nil {
				c.HTML(http.StatusOK, "error.gohtml", gin.H{
					"error": e,
				})
				log.Panicln("無法保存文件", e.Error())
			}
			avatarUrl := "http://localhost:8080/avatar/" + fileName
			requestBody.Avatar = sql.NullString{String: avatarUrl}
		}

		item := user.Item{
			ID:       requestBody.Id,
			Email:    requestBody.Email,
			Password: requestBody.Password,
			Avatar:   requestBody.Avatar,
			Name:     requestBody.Name,
		}

		if item.Avatar.String == "" {
			e = userUpdater.UpdateNoAvatar(item)
			if e != nil {
				c.HTML(http.StatusOK, "error.gohtml", gin.H{
					"error": e,
				})
				log.Panicln("数据无法更新", e.Error())
			}
		} else {
			e = userUpdater.Update(item)
			if e != nil {
				c.HTML(http.StatusOK, "error.gohtml", gin.H{
					"error": e,
				})
				log.Panicln("数据无法更新", e.Error())
			}
		}

		if claims := middleware.GetCurentUserClaims(c); claims != nil {
			c.HTML(http.StatusOK, "index.gohtml", gin.H{
				"result": result,
				"email":  claims.Audience,
				"id":     string(claims.Id),
			})
		}
	}
}
