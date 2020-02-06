package articlehandler

import (
	"inventory/httpd/model"
	"inventory/platform/article"
	"net/http"

	"github.com/gin-gonic/gin"
)

type articleAddPostResquest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// @Summary 新增文章
// @Id 1
// @Tags 文章
// @version 1.0
// @Accept application/x-json-stream
// @Param article body articlehandler.articleAddPostResquest true "文章"
// @Success 200 object model.Result 成功回傳
// @Failure 409 object model.Result 失敗回傳
// @Router /article [post]
func ArticleAdd(articleAdder article.Adder) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := articleAddPostResquest{}
		var id = -1
		var message = "新增失敗"
		if e := c.ShouldBindJSON(&requestBody); e == nil {
			item := article.Item{
				Type:    requestBody.Type,
				Content: requestBody.Content,
			}
			id = articleAdder.Add(item)
			message = "新增成功"
		}

		result := model.Result{
			Code:    http.StatusOK,
			Message: message,
			Data: gin.H{
				"id": id,
			},
		}
		c.JSON(http.StatusOK, result)
	}
}
