package userhandler

import (
	"inventory/httpd/middleware"
	"inventory/platform/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserProfileGet(userGetter user.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if claims := middleware.GetCurentUserClaims(c); claims != nil {
			if claims.Id != id { //查別人的資料
				c.HTML(http.StatusUnauthorized, "401.gohtml", gin.H{
					"error": "不要亂來哦",
				})
				c.Abort()
			} else {
				i, err := strconv.Atoi(id)
				u, e := userGetter.QueryById(i)
				if e != nil || err != nil {
					c.HTML(http.StatusOK, "error.gohtml", gin.H{
						"error": e,
					})
				}
				c.HTML(http.StatusOK, "user_profile.gohtml", gin.H{
					"email": claims.Audience,
					"id":    string(claims.Id),
					"user":  u,
				})
			}
		} else {
			c.HTML(http.StatusUnauthorized, "401.gohtml", nil)
			c.Abort()
		}
		//c.JSON(http.StatusOK, u)
	}
}
