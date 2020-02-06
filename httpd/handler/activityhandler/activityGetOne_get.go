package activityhandler

import (
	"inventory/httpd/middleware"
	"inventory/platform/activity"
	"inventory/platform/activityApply"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ActivityGetOne(activityGetter activity.Getter, activityApplyAction activityApply.Action) gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")
		i, _ := strconv.Atoi(id)

		activity, _ := activityGetter.FindById(i)
		if claims := middleware.GetCurentUserClaims(c); claims != nil {

			applyList := activityApplyAction.GetApplyList(i)
			c.HTML(http.StatusOK, "activity_detail.gohtml", gin.H{
				"email":     claims.Audience,
				"id":        claims.Id,
				"activity":  activity,
				"applyList": applyList,
				"page":      "list",
			})
		} else {
			c.HTML(http.StatusOK, "401.gohtml", gin.H{})
		}

	}
}
