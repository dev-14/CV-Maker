package routes

import (
	//_ "gingorm/docs"
	// middlewares "gingorm/middleware"
	"CV-Maker/controllers"
	middlewares "CV-Maker/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"
)

func PublicEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// Generate public endpoints - [ signup] - api/v1/signup

	r.POST("/register", controllers.Register)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/logout", authMiddleware.LogoutHandler)
	r.GET("/refresh", authMiddleware.RefreshHandler)
}

func AuthenticatedEndpoints(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	r.Use(authMiddleware.MiddlewareFunc())

	// education endpoints
	r.POST("/education", controllers.CreateEducation)
	r.GET("education/:id", controllers.GetEducation)
	r.GET("/education", controllers.GetAllEducation)
	r.GET("/education?year", controllers.GetEducationByYear)
	r.PATCH("/education/:id", controllers.UpdateEducation)
	r.DELETE("/education/:id", controllers.DeleteEducation)

	// workExperience endpoints

	r.POST("/work", controllers.CreateWorkExperience)
	r.GET("/work/:id", controllers.GetWorkExperience)
	r.GET("/work", controllers.GetAllWorkExperience)
	r.GET("/work/getByYear", controllers.GetWorkExperienceByYear)
	r.PATCH("/work/:id", controllers.UpdateWorkExperience)
	r.DELETE("/work/:id", controllers.DeleteWorkExperience)

	// project endpoints

	r.POST("/project", controllers.CreateProject)
	r.GET("project/:id", controllers.GetProject)
	r.GET("/project", controllers.GetAllProject)
	r.PATCH("/project/:id", controllers.UpdateProject)
	r.DELETE("/project/:id", controllers.DeleteProject)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Headers", "Content-Type, redirect, body, Content-Length, Authentication, Accept-Encoding, X-CSRF-Token, Authorization, method, accept, origin, Cache-Control, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func GetRouter(router chan *gin.Engine) {
	gin.ForceConsoleColor()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	_ = config
	r.Use(CORSMiddleware())
	// r.Use(middlewares.RequestLogger)
	// r.Use(gin.CustomRecovery(middlewares.LogFailedRequests))
	authMiddleware, _ := middlewares.GetAuthMiddleware()

	// Create a BASE_URL - /api/v1
	v1 := r.Group("/api/v1/")
	PublicEndpoints(v1, authMiddleware)
	AuthenticatedEndpoints(v1.Group("auth"), authMiddleware)
	router <- r
}
