package main

import (
	"blog-system/internal/config"
	"blog-system/internal/controller"
	"blog-system/internal/middleware"
	"blog-system/internal/repository"
	"blog-system/internal/service"
	"blog-system/internal/utils"
	"blog-system/pkg/logger"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志
	appLogger := logger.NewLogger()

	// 初始化数据库
	db, err := utils.InitDB(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化JWT工具
	jwtUtil := utils.NewJWTUtil(cfg)

	// 初始化Repository
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// 初始化Service
	authService := service.NewAuthService(userRepo, jwtUtil)
	postService := service.NewPostService(postRepo, userRepo, appLogger)
	commentService := service.NewCommentService(commentRepo, postRepo, userRepo, appLogger)

	// 初始化Controller
	authController := controller.NewAuthController(authService)
	postController := controller.NewPostController(postService)
	commentController := controller.NewCommentController(commentService)

	// 创建Gin应用
	app := setupRouter(authController, postController, commentController, jwtUtil, appLogger)

	// 启动服务器
	appLogger.Info("服务器启动在端口: %s", cfg.ServerPort)
	if err := app.Run(":" + cfg.ServerPort); err != nil {
		appLogger.Fatal("服务器启动失败: %v", err)
	}
}

func setupRouter(
	authController *controller.AuthController,
	postController *controller.PostController,
	commentController *controller.CommentController,
	jwtUtil *utils.JWTUtil,
	appLogger *logger.Logger,
) *gin.Engine {
	app := gin.Default()

	// 中间件
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.LoggingMiddleware(appLogger))

	// 健康检查
	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog System API is running",
		})
	})

	// API路由组
	api := app.Group("/api")

	// 认证路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/profile", middleware.AuthMiddleware(jwtUtil), authController.GetProfile)
	}

	// 文章路由
	posts := api.Group("/posts")
	{
		posts.GET("", postController.GetPosts)
		posts.GET("/:id", postController.GetPost)

		// 需要认证的路由
		authPosts := posts.Group("", middleware.AuthMiddleware(jwtUtil))
		{
			authPosts.POST("", postController.CreatePost)
			authPosts.PUT("/:id", postController.UpdatePost)
			authPosts.DELETE("/:id", postController.DeletePost)
			authPosts.GET("/my", postController.GetUserPosts)      // 获取所有文章列表
			authPosts.GET("/user/:id", postController.GetUserPost) // 获取单个文章的详细信息
		}
	}

	// 评论路由
	comments := api.Group("/comments")
	{
		comments.GET("/posts/:postId", commentController.GetCommentsByPost)

		// 需要认证的路由
		authComments := comments.Group("", middleware.AuthMiddleware(jwtUtil))
		{
			authComments.POST("", commentController.CreateComment)
			authComments.GET("/my", commentController.GetUserComments)
			authComments.DELETE("/:id", commentController.DeleteComment)
		}
	}

	// 404处理
	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"success": false,
			"message": "接口不存在",
			"data":    nil,
		})
	})

	return app
}

// 初始化数据库表和数据
func initDatabase(db *gorm.DB, appLogger *logger.Logger) {
	// 可以在这里添加初始数据
	appLogger.Info("数据库初始化完成")
}
