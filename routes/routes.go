package routes

import (
	"net/http"
	"time"

	"github.com/disturb-yy/bluebell/middlewares"

	"github.com/disturb-yy/bluebell/controller"

	"github.com/disturb-yy/bluebell/settings"

	"github.com/disturb-yy/bluebell/logger"

	_ "github.com/disturb-yy/bluebell/docs" // 千万不要忘了导入把你上一步生成的docs

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/gin-contrib/pprof"
	// swagger embed files
	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		// 发布模式，否则默认为开发模式
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	// 注册中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 加载静态文件
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 路由匹配
	v1 := r.Group("/api/v1")
	{
		// 注册
		v1.POST("/signup", controller.SignUpHandler)
		// 登录
		v1.POST("/login", controller.LoginHandler)

		// 获取帖子list显示，posts按创建排序，posts2按分数排序
		v1.GET("/posts", controller.GetPostListHandler)
		// 可以携带社区id
		v1.GET("/posts2", controller.GetPostListHandler2)
		// 根据 url 传入的 id，获取对应的帖子详情
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		// 在路由组v1里使用中间件
		//验证模块和令牌限流
		v1.Use(middlewares.JWTAuthMiddleware(), middlewares.RateLimitMiddleware(2*time.Second, 1)) // 应用JWT认证中间件
		{
			v1.POST("/post", controller.CreatePostHandler)
			// 实现用户投票功能
			v1.POST("/vote", controller.PostVoteHandler)
		}
	}

	// 响应
	r.GET("/ping", func(c *gin.Context) {
		// 如果是登录的用户，判断请求头中是否有 有效的JWT
		c.String(http.StatusOK, "pong")
	})

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})

	pprof.Register(r) // 注册 pprof 相关路由
	//r.Run()
	return r
}
