package middleware

import (
	"inventory/config"
	"inventory/httpd/initDB"
	"inventory/platform/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// func Auth() gin.HandlerFunc {
// 	return func(context *gin.Context) {
// 		cookie, e := context.Request.Cookie("user_cookie")
// 		if e == nil {
// 			context.SetCookie(cookie.Name, cookie.Value, 1000, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
// 			context.Next()
// 		} else {
// 			context.Abort()
// 			context.HTML(http.StatusUnauthorized, "401.gohtml", nil)
// 		}
// 	}
// }
// func Auth() gin.HandlerFunc {
// 	return func(context *gin.Context) {
// 		result := model.Result{
// 			Code:    http.StatusUnauthorized,
// 			Message: "無法認證，請重新登入",
// 			Data:    nil,
// 		}
// 		auth := context.Request.Header.Get("Authorization")
// 		if len(auth) == 0 {
// 			context.Abort()
// 			context.JSON(http.StatusUnauthorized, gin.H{
// 				"result": result,
// 			})
// 		}
// 		auth = strings.Fields(auth)[1]
// 		// 驗證token
// 		_, err := parseToken(auth)
// 		if err != nil {
// 			context.Abort()
// 			result.Message = "token 過期" + err.Error()
// 			context.JSON(http.StatusUnauthorized, gin.H{
// 				"result": result,
// 			})
// 		} else {
// 			println("token 正確")
// 		}
// 		context.Next()
// 	}
// }

func GetCurentUserClaims(context *gin.Context) *jwt.StandardClaims {

	if claims, _ := context.Get("user"); claims != nil {
		if claims, ok := claims.(*jwt.StandardClaims); ok {
			return claims
		}
	}
	return nil
}

func Auth(failHtml string) gin.HandlerFunc {
	return func(context *gin.Context) {
		authCookie, e := context.Request.Cookie("user_auth")
		redirectPage := failHtml + ".gohtml"
		if e == nil {
			//auth := strings.Fields(authCookie.Value)[1]
			auth := authCookie.Value
			// 驗證token
			claim, err := parseToken(auth)
			if err != nil {
				if failHtml != "401" {
					context.Next()
				} else {
					context.Abort()
					context.HTML(http.StatusUnauthorized, redirectPage, nil)
				}
			} else {
				println("token 正確")
				context.Set("user", claim)
				context.SetCookie(authCookie.Name, authCookie.Value, config.CookieAliveTime, authCookie.Path, authCookie.Domain, authCookie.Secure, authCookie.HttpOnly)
				context.Next()
			}
		} else {
			if failHtml != "401" {
				context.Next()
			} else {
				context.Abort()
				context.HTML(http.StatusUnauthorized, redirectPage, nil)
			}
		}
	}
}

func parseToken(token string) (*jwt.StandardClaims, error) {
	// 分割出來 payload
	payload := strings.Split(token, ".")
	bytes, e := jwt.DecodeSegment(payload[1])

	if e != nil {
		println(e.Error())
	}
	content := ""
	for i := 0; i < len(bytes); i++ {
		content += string(bytes[i])
	}
	split := strings.Split(content, ",")
	id := strings.SplitAfter(split[2], ":")
	ID, err := strconv.Atoi(strings.ReplaceAll(id[1], "\"", ""))
	if err != nil {
		println(err.Error())
	}

	Db := initDB.Db
	user, _ := user.Getter(user.NewUser(Db)).QueryById(ID)

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(config.Secret + user.Password), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}
