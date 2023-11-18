package main

import (
	"identifEye/database"
	"identifEye/middleware"
	"identifEye/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load env file: ", err)
	}

	database.Init()

	r := gin.Default()
	r.Use(gin.Logger())

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
