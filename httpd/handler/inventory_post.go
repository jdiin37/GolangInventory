package handler

import (
	"inventory/platform/inventory"
	"net/http"

	"github.com/gin-gonic/gin"
)

type inventoryPostResquest struct {
	Title string `json:"title"`
	Post  string `json:"post"`
}

func InventoryPost(inventorys inventory.Adder) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := inventoryPostResquest{}
		c.Bind(&requestBody)

		item := inventory.Item{
			Title: requestBody.Title,
			Post:  requestBody.Post,
		}

		inventorys.Add(item)
		c.Status(http.StatusNoContent)
	}
}
