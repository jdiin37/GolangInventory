package articlehandler

import (
	"inventory/httpd/model"
	"inventory/platform/article"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 取得文章
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Param id path int true "id"
// @Success 200 {object} model.Result 成功回傳
// @Router /article/{id} [get]
func ArticleGetOne(articleGetter article.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		i, e := strconv.Atoi(id)
		if e != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": model.Result{
					Code:    http.StatusBadRequest,
					Message: "id 不是 int 类型, id 转换失败",
					Data:    e.Error(),
				},
			})
			log.Panicln("id 不是 int 类型, id 转换失败", e.Error())
		}

		if article, err := articleGetter.FindById(i); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": model.Result{
					Code:    http.StatusBadRequest,
					Message: "get error",
					Data:    err.Error(),
				},
			})
			log.Panicln("get error", e.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": model.Result{
					Code:    http.StatusOK,
					Message: "查詢成功",
					Data:    article,
				},
			})
		}

	}
}
