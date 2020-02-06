package main

import (
	_ "inventory/httpd/docs"
	"inventory/httpd/initRouter"
	"net/http"
	"time"
)

// @title Gin swagger
// @version 1.0
// @description Gin swagger

// @contact.name jdiin37
// @contact.email jdiin37@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080/api/v1
func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := initRouter.SetupRouter()

	//router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")d
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
	// routerHTTPS := initRouter.SetupRouter()
	// routerHTTPS.RunTLS(":8888", "./server.crt", "./server.key")
}
