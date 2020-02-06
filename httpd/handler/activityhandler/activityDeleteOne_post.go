package activityhandler

import (
	"inventory/httpd/model"
	"inventory/platform/activity"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type activityDeleteOnePostResquest struct {
	ID int `form:"activityId"`
}

func ActivityDeleteOne(activityDeleter activity.Deleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := activityDeleteOnePostResquest{}
		result := &model.Result{
			Code:    http.StatusBadRequest,
			Message: "刪除失敗",
			Data:    nil,
		}
		if e := c.ShouldBind(&requestBody); e == nil {
			activityDeleter.DeleteOne(requestBody.ID)
			result.Code = http.StatusOK
			result.Message = "刪除成功"
			result.Data = gin.H{
				"activityID": requestBody.ID,
			}
			//c.JSON(http.StatusOK, result)
			c.Redirect(302, "/activity/list")
		} else {
			log.Println(e.Error())
			c.Redirect(400, "/error")
		}
	}
}
