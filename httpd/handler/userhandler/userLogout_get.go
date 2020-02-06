package userhandler

import (
	"inventory/httpd/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserLogoutGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		result := &model.Result{
			Code:    200,
			Message: "登出成功",
			Data:    nil,
		}
		c.SetCookie("user_auth", "", -1, "/", "localhost", false, true)
		c.HTML(http.StatusOK, "index.gohtml", gin.H{
			"result": result,
		})
		//c.JSON(http.StatusOK, u)
	}
}
