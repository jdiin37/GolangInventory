package articlehandler

import (
	"inventory/httpd/model"
	"inventory/platform/article"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 取得所有文章
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Success 200 object model.Result 成功回傳
// @Router /articles [get]
func ArticleGetAll(articleGetter article.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			offset = 0
			log.Println("====== offset nil ======")
		}
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 1000
			log.Println("====== limit nil ======")
		}
		articles := articleGetter.GetAll(offset, limit)

		result := model.Result{
			Code:    http.StatusOK,
			Message: "查詢成功",
			Data:    articles,
		}
		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}
