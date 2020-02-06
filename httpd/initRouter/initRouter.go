package initRouter

import (
	"html/template"
	"inventory/httpd/handler"
	"inventory/httpd/handler/activityhandler"
	"inventory/httpd/handler/articlehandler"
	"inventory/httpd/handler/userhandler"
	"inventory/httpd/initDB"
	"inventory/httpd/middleware"
	"inventory/kit"
	"inventory/platform/activity"
	"inventory/platform/activityApply"
	"inventory/platform/article"
	"inventory/platform/inventory"
	"inventory/platform/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	// Add whatever other types you need
	default:
		return ""
	}
}

func add(left int, right int) int {
	return left + right
}

func SetupRouter() *gin.Engine {

	Db := initDB.Db

	r := gin.New()
	//r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"toString": ToString,
		"add":      add,
	})

	if mode := gin.Mode(); mode == gin.TestMode {
		r.LoadHTMLGlob("./../templates/*")
	} else {
		//r.Delims("{[{", "}]}")
		r.LoadHTMLGlob("./httpd/templates/*")
	}

	r.Use(middleware.Logger(), gin.Recovery())
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/inventory", handler.InventoryGet(inventory.NewInventory(Db)))
			v1.POST("/inventory", middleware.Auth("401"), handler.InventoryPost(inventory.NewInventory(Db)))

			v1.POST("/userRegister", userhandler.UserRegisterPost(user.NewUser(Db)))
			v1.POST("/userLogin", userhandler.UserLoginPost(user.NewUser(Db)))
			v1.POST("/userProfile", middleware.Auth("401"), userhandler.UserProfilePost(user.NewUser(Db)))

			// 取得一篇文章
			v1.GET("/article/:id", articlehandler.ArticleGetOne(article.NewArticle(Db)))
			// 取得所有文章
			v1.GET("/articles", articlehandler.ArticleGetAll(article.NewArticle(Db)))
			// 添加一篇文章
			v1.POST("/article", articlehandler.ArticleAdd(article.NewArticle(Db)))
			v1.DELETE("/article/:id", articlehandler.ArticleDeleteOne(article.NewArticle(Db)))

			//參加團
			v1.POST("activityApply", middleware.Auth("activity"), activityhandler.ActivityApply(activityApply.NewActivityApply(Db)))

		}
	}
	r.Static("/statics", "./httpd/statics")
	r.StaticFile("/favicon.ico", "./httpd/statics/favicon.ico")
	r.StaticFS("/avatar", http.Dir(kit.RootPath()+"avatar"))
	index := r.Group("/")
	{
		index.GET("", middleware.Auth("index"), handler.IndexGet())
	}

	userRouter := r.Group("/user")
	{
		userRouter.GET("/logout", userhandler.UserLogoutGet())
		userRouter.GET("/profile", middleware.Auth("401"), userhandler.UserProfileGet(user.NewUser(Db)))
	}

	activityRouter := r.Group("/activity")
	{
		activityRouter.GET("/list", middleware.Auth("activity"), activityhandler.ActivityGetAll(activity.NewActivity(Db)))
		activityRouter.POST("/", middleware.Auth("activity"), activityhandler.ActivityAdd(activity.NewActivity(Db)))
		activityRouter.POST("/DeleteOne", middleware.Auth("activity"), activityhandler.ActivityDeleteOne(activity.NewActivity(Db)))
		activityRouter.GET("/detail/:id", middleware.Auth("activity"),
			activityhandler.ActivityGetOne(activity.NewActivity(Db), activityApply.NewActivityApply(Db)))
	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}
