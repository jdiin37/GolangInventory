package activityhandler

import (
	"inventory/httpd/middleware"
	"inventory/httpd/model"
	"inventory/platform/activityApply"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type activityApplyPostResquest struct {
	ActivityID int    `json:"activityID"`
	Remark     string `json:"remark"`
}

func ActivityApply(activityAction activityApply.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := &model.Result{
			Code:    http.StatusBadRequest,
			Message: "ActivityApply 失敗",
			Data:    nil,
		}
		if claims := middleware.GetCurentUserClaims(c); claims != nil {
			requestBody := activityApplyPostResquest{}
			userID, _ := strconv.Atoi(claims.Id)

			if e := c.ShouldBindJSON(&requestBody); e == nil {
				item := activityApply.Item{
					Remark:     requestBody.Remark,
					ActivityID: requestBody.ActivityID,
					UserID:     userID,
				}
				if activityAction.CheckCanApply(item) {
					applyID := activityAction.UserApply(item)
					result.Code = http.StatusOK
					result.Message = "ActivityApply 成功"
					result.Data = gin.H{
						"applyID": applyID,
					}
				} else {
					result.Code = http.StatusBadRequest
					result.Message = "ActivityApply 已報名"
					result.Data = gin.H{}
				}

			} else {
				log.Println(e.Error())
			}
		}
		c.JSON(http.StatusOK, result)
	}
}
