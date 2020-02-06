package activityhandler

import (
	"inventory/httpd/middleware"
	"inventory/httpd/model"
	"inventory/platform/activityApply"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ActivityApplyGetList(activityAction activityApply.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := &model.Result{
			Code:    http.StatusBadRequest,
			Message: "ActivityApplyGetList 失敗",
			Data:    nil,
		}
		activityID := c.Param("activityID")
		id, _ := strconv.Atoi(activityID)
		if claims := middleware.GetCurentUserClaims(c); claims != nil {
			data := activityAction.GetApplyList(id)
			result.Code = http.StatusOK
			result.Message = "ActivityApplyGetList 成功"
			result.Data = gin.H{
				"data": data,
			}
		} else {
			result.Message = "ActivityApplyGetList 失敗 未登入"
		}
		c.JSON(http.StatusOK, result)
	}
}
