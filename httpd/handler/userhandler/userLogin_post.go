package userhandler

import (
	"inventory/config"
	"inventory/httpd/model"
	"inventory/platform/user"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type userLoginPostResquest struct {
	Email    string `form:"email" binding:"email"`
	Password string `form:"password"`
}

func UserLoginPost(users user.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := userLoginPostResquest{}
		result := &model.Result{
			Code:    200,
			Message: "登入成功",
			Data:    nil,
		}
		if e := c.Bind(&requestBody); e != nil {
			//log.Panicln("login 綁定錯誤", e.Error())
			result.Message = "資料格式綁定錯誤"
			result.Code = http.StatusBadRequest
			// c.JSON(http.StatusUnauthorized, gin.H{
			// 	"result": result,
			// })
			c.HTML(http.StatusOK, "index.gohtml", gin.H{
				"result": result,
			})
			c.Abort()
			return
		}

		u := users.QueryByEmail(requestBody.Email)
		if u.Password == requestBody.Password {
			log.Println("登入成功", u.Email)
			//c.SetCookie("user_cookie", string(u.ID), 1000, "/", "localhost", false, true)
			// c.HTML(http.StatusOK, "index.gohtml", gin.H{
			// 	"email": u.Email,
			// 	"id":    u.ID,
			// })
			expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
			claims := jwt.StandardClaims{
				Audience:  u.Email,            // 受众
				ExpiresAt: expiresTime,        // 失效时间
				Id:        strconv.Itoa(u.ID), // 编号
				IssuedAt:  time.Now().Unix(),  // 签发时间
				Issuer:    "jdiin",            // 签发人
				NotBefore: time.Now().Unix(),  // 生效时间
				Subject:   "login",            // 主题
			}

			tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// 透過密碼和保留字段加密
			var jwtSecret = []byte(config.Secret + u.Password)
			if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
				result.Message = "登入成功"
				result.Data = "Bearer " + token
				result.Code = http.StatusOK
				// c.JSON(result.Code, gin.H{
				// 	"result": result,
				// })
				c.SetCookie("user_auth", string(token), config.CookieAliveTime, "/", "localhost", false, true)
				c.HTML(http.StatusOK, "index.gohtml", gin.H{
					"result": result,
					"email":  u.Email,
					"id":     u.ID,
				})
			} else {
				result.Message = "登入失敗"
				result.Code = http.StatusBadRequest
				// c.JSON(result.Code, gin.H{
				// 	"result": result,
				// })
				c.HTML(http.StatusOK, "index.gohtml", gin.H{
					"result": result,
				})
			}
		} else {
			result.Message = "登入失敗"
			result.Code = http.StatusBadRequest
			// c.JSON(http.StatusUnauthorized, gin.H{
			// 	"result": result,
			// })
			c.HTML(http.StatusOK, "index.gohtml", gin.H{
				"result": result,
			})
		}

	}
}
