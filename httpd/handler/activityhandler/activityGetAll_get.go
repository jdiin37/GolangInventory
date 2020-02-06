package activityhandler

import (
	"inventory/httpd/middleware"
	"inventory/platform/activity"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ActivityGetAll(activityGetter activity.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {

		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			offset = 0
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}

		activitys := activityGetter.GetAll(offset, limit)

		if claims := middleware.GetCurentUserClaims(c); claims != nil {
			c.HTML(http.StatusOK, "activity.gohtml", gin.H{
				"email":     claims.Audience,
				"id":        claims.Id,
				"activitys": activitys,
				"page":      "list",
			})
		} else {
			c.HTML(http.StatusOK, "activity.gohtml", gin.H{
				"activitys": activitys,
				"page":      "list",
			})
		}

	}
}
