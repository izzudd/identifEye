package main

import (
	"identifEye/database"
	"identifEye/middleware"
	"identifEye/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load env file: ", err)
	}

	database.Init()

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	r.Use(gin.Logger())
	r.Use(cors.New(config))

	r.POST("/register", routes.RegisterHandler)
	r.POST("/register/face", middleware.AuthorizeRegister, routes.RegisterFaceHandler)
	r.POST("/login", routes.LoginHandler)
	r.POST("/login/face", middleware.AuthorizeNoFace, routes.LoginFaceHandler)
	r.GET("/secret", middleware.AuthorizeNoFace, middleware.AuthorizeFace, routes.Secret)

	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
