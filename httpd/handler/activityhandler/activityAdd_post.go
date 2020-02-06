package activityhandler

import (
	"inventory/httpd/middleware"
	"inventory/httpd/model"
	"inventory/platform/activity"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type activityAddPostResquest struct {
	StartDate string `form:"startDate"`
	StartTime string `form:"startTime"`
	EndDate   string `form:"endDate"`
	EndTime   string `form:"endTime"`
	County    string `form:"county"`
	Location  string `form:"location"`
	CourtType string `form:"courtType"`
	Memo      string `form:"memo"`
	Context   string `form:"context"`
}

func ActivityAdd(activityAdder activity.Adder) gin.HandlerFunc {
	return func(c *gin.Context) {

		if claims := middleware.GetCurentUserClaims(c); claims != nil {
			requestBody := activityAddPostResquest{}
			userID, _ := strconv.Atoi(claims.Id)
			var activityID = -1
			result := &model.Result{
				Code:    http.StatusBadRequest,
				Message: "新增失敗",
				Data:    nil,
			}
			if e := c.ShouldBind(&requestBody); e == nil {
				item := activity.Item{
					StartTime:    requestBody.StartDate + " " + requestBody.StartTime,
					EndTime:      requestBody.EndDate + " " + requestBody.EndTime,
					County:       requestBody.County,
					Location:     requestBody.Location,
					CourtType:    requestBody.CourtType,
					Memo:         requestBody.Memo,
					Context:      requestBody.Context,
					CreateUserID: userID,
				}
				activityID = activityAdder.Add(item)
				result.Code = http.StatusOK
				result.Message = "新增成功"
				result.Data = gin.H{
					"activityID": activityID,
				}
				//c.JSON(http.StatusOK, result)
				c.Redirect(302, "/activity/list")
			} else {
				log.Println(e.Error())
				c.Redirect(400, "/error")
			}
			c.Redirect(302, "/activity/list")
			// c.HTML(http.StatusOK, "index.gohtml", gin.H{
			// 	"result": result,
			// 	"email":  claims.Audience,
			// 	"id":     string(claims.Id),
			// })
		} else {
			c.Abort()
		}
	}
}
