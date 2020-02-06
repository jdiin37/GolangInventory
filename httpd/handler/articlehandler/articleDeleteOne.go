package articlehandler

import (
	"inventory/platform/article"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 刪除文章
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Param id path int true "id"
// @Success 200  object model.Result 成功回傳
// @Router /article/{id} [delete]
func ArticleDeleteOne(articleDeleter article.Deleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		i, e := strconv.Atoi(id)
		if e != nil {
			log.Panicln("id 不是 int 类型, id 转换失败", e.Error())
		}
		articleDeleter.DeleteOne(i)
	}
}
