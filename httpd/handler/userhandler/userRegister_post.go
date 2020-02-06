package userhandler

import (
	"inventory/httpd/model"
	"inventory/platform/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userRegisterPostResquest struct {
	Email         string `form:"email" binding:"email"`
	Password      string `form:"password"`
	PasswordAgain string `form:"password-again" binding:"eqfield=Password"`
}

func UserRegisterPost(users user.Adder) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := userRegisterPostResquest{}

		result := &model.Result{
			Code:    200,
			Message: "登入成功",
			Data:    nil,
		}
		if err := c.ShouldBind(&requestBody); err != nil {
			result.Message = "輸入的數據不合法"
			result.Code = http.StatusUnauthorized
			// c.JSON(http.StatusUnauthorized, gin.H{
			// 	"result": result,
			// })
			c.HTML(http.StatusOK, "index.gohtml", gin.H{
				"result": result,
			})
		} else {
			item := user.Item{
				Email:    requestBody.Email,
				Password: requestBody.Password,
			}
			if err := users.Add(item); err != nil {
				result.Message = "註冊失败"
				result.Code = http.StatusOK

				c.JSON(http.StatusOK, gin.H{
					"result": result,
				})
			} else {
				result.Message = "註冊成功"
				result.Code = http.StatusOK
				c.JSON(http.StatusOK, gin.H{
					"result": result,
				})
			}
			//c.Status(http.StatusNoContent)
			// log.Println("email", requestBody.Email, "password", requestBody.Password, "password again", requestBody.PasswordAgain)
			// c.Redirect(http.StatusMovedPermanently, "/")
		}

	}
}
