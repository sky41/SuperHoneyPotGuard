package routes

import (
	"superhoneypotguard/controllers"
	"superhoneypotguard/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	roleController := controllers.NewRoleController()
	permissionController := controllers.NewPermissionController()
	dashboardController := controllers.NewDashboardController()
	logController := controllers.NewLogController()
	hfishController := controllers.NewHFishController()
	passwordController := controllers.NewPasswordController()

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success":   true,
				"message":   "服务运行正常",
				"timestamp": gin.H{},
			})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/send-verification-code", authController.SendVerificationCode)
			auth.POST("/send-reset-code", passwordController.SendResetPasswordCode)
			auth.POST("/register", authController.Register)
			auth.POST("/login", middleware.AuthRateLimitMiddleware(), authController.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), authController.Logout)
			auth.GET("/current", middleware.AuthMiddleware(), authController.GetCurrentUser)
		}

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/list", middleware.PermissionMiddleware("user:manage"), userController.GetList)
			user.GET("/:id", middleware.PermissionMiddleware("user:manage"), userController.GetById)
			user.POST("/", middleware.PermissionMiddleware("user:manage"), userController.Create)
			user.PUT("/:id", middleware.PermissionMiddleware("user:manage"), userController.Update)
			user.DELETE("/:id", middleware.PermissionMiddleware("user:manage"), userController.Delete)
			user.PATCH("/:id/status", middleware.PermissionMiddleware("user:manage"), userController.UpdateStatus)
			user.POST("/:id/reset-password", middleware.PermissionMiddleware("user:manage"), userController.ResetPassword)
		}

		role := api.Group("/role")
		role.Use(middleware.AuthMiddleware())
		{
			role.GET("/list", middleware.PermissionMiddleware("role:manage"), roleController.GetList)
			role.GET("/all", middleware.PermissionMiddleware("role:manage"), roleController.GetAll)
			role.GET("/:id", middleware.PermissionMiddleware("role:manage"), roleController.GetById)
			role.POST("/", middleware.PermissionMiddleware("role:manage"), roleController.Create)
			role.PUT("/:id", middleware.PermissionMiddleware("role:manage"), roleController.Update)
			role.DELETE("/:id", middleware.PermissionMiddleware("role:manage"), roleController.Delete)
		}

		permission := api.Group("/permission")
		permission.Use(middleware.AuthMiddleware())
		{
			permission.GET("/tree", middleware.PermissionMiddleware("permission:manage"), permissionController.GetTree)
			permission.GET("/all", middleware.PermissionMiddleware("permission:manage"), permissionController.GetAll)
			permission.GET("/:id", middleware.PermissionMiddleware("permission:manage"), permissionController.GetById)
			permission.POST("/", middleware.PermissionMiddleware("permission:manage"), permissionController.Create)
			permission.PUT("/:id", middleware.PermissionMiddleware("permission:manage"), permissionController.Update)
			permission.DELETE("/:id", middleware.PermissionMiddleware("permission:manage"), permissionController.Delete)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.AuthMiddleware())
		{
			dashboard.GET("/stats", middleware.PermissionMiddleware("dashboard:view"), dashboardController.GetStats)
		}

		log := api.Group("/log")
		log.Use(middleware.AuthMiddleware())
		{
			log.GET("/list", middleware.PermissionMiddleware("log:manage"), logController.GetList)
			log.GET("/:id", middleware.PermissionMiddleware("log:manage"), logController.GetById)
			log.DELETE("/:id", middleware.PermissionMiddleware("log:manage"), logController.Delete)
			log.DELETE("/clear", middleware.PermissionMiddleware("log:manage"), logController.Clear)
		}

		hfish := api.Group("/hfish")
		hfish.Use(middleware.AuthMiddleware())
		{
			hfish.GET("/attack/ips", middleware.PermissionMiddleware("hfish:view"), hfishController.GetAttackIPs)
			hfish.GET("/attack/details", middleware.PermissionMiddleware("hfish:view"), hfishController.GetAttackDetails)
			hfish.GET("/account/info", middleware.PermissionMiddleware("hfish:view"), hfishController.GetAccountInfo)
			hfish.GET("/sys/info", middleware.PermissionMiddleware("hfish:view"), hfishController.GetSysInfo)
			hfish.POST("/block/ip", middleware.PermissionMiddleware("hfish:block"), hfishController.BlockIP)
		}

		password := api.Group("/password")
		password.Use(middleware.AuthMiddleware())
		{
			password.POST("/send-reset-code", passwordController.SendResetPasswordCode)
			password.POST("/reset", passwordController.ResetPassword)
		}

		r.NoRoute(func(c *gin.Context) {
			c.JSON(404, gin.H{
				"success": false,
				"message": "请求的资源不存在",
			})
		})
	}
}
