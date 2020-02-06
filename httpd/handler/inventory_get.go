package handler

import (
	"inventory/platform/inventory"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type reqQueryString struct {
	offset string `form:"offset"`
	limit  string `form:"limit"`
}

func InventoryGet(i interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch obj := i.(type) {
		case inventory.Getter:

			// queryString := reqQueryString{}

			// c.Bind(&queryString)
			offset, err := strconv.Atoi(c.Query("offset"))
			if err != nil {
				offset = 0
				log.Println("====== offset nil ======")
			}
			limit, err := strconv.Atoi(c.Query("limit"))
			if err != nil {
				limit = 10
				log.Println("====== limit nil ======")
			}

			results := obj.Get(offset, limit)
			c.JSON(http.StatusOK, results)

		}

		// _, isGet := inventory.(Getter)

		// if isGet {
		// 	results := inventory.Get()
		// 	c.JSON(http.StatusOK, results)
		// }
	}
}
