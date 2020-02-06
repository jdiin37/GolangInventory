package handler

import (
	"inventory/httpd/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims := middleware.GetCurentUserClaims(c); claims != nil {
			c.HTML(http.StatusOK, "index.gohtml", gin.H{
				"email": claims.Audience,
				"id":    string(claims.Id),
			})
		} else {
			c.HTML(http.StatusOK, "index.gohtml", gin.H{})
		}

	}
}
