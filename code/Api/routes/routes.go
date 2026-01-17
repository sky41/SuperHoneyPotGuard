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
			auth.POST("/register", authController.Register)
			auth.POST("/login", middleware.AuthRateLimitMiddleware(), authController.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), authController.Logout)
			auth.GET("/current", middleware.AuthMiddleware(), authController.GetCurrentUser)
		}

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/list", middleware.PermissionMiddleware("user:list"), userController.GetList)
			user.GET("/:id", middleware.PermissionMiddleware("user:list"), userController.GetById)
			user.POST("/", middleware.PermissionMiddleware("user:add"), userController.Create)
			user.PUT("/:id", middleware.PermissionMiddleware("user:edit"), userController.Update)
			user.DELETE("/:id", middleware.PermissionMiddleware("user:delete"), userController.Delete)
			user.PATCH("/:id/status", middleware.PermissionMiddleware("user:edit"), userController.UpdateStatus)
			user.POST("/:id/reset-password", middleware.PermissionMiddleware("user:edit"), userController.ResetPassword)
		}

		role := api.Group("/role")
		role.Use(middleware.AuthMiddleware())
		{
			role.GET("/list", middleware.PermissionMiddleware("role:list"), roleController.GetList)
			role.GET("/all", middleware.PermissionMiddleware("role:list"), roleController.GetAll)
			role.GET("/:id", middleware.PermissionMiddleware("role:list"), roleController.GetById)
			role.POST("/", middleware.PermissionMiddleware("role:add"), roleController.Create)
			role.PUT("/:id", middleware.PermissionMiddleware("role:edit"), roleController.Update)
			role.DELETE("/:id", middleware.PermissionMiddleware("role:delete"), roleController.Delete)
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
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"message": "请求的资源不存在",
		})
	})
}
