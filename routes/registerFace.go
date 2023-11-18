package routes

import (
	"fmt"
	"identifEye/database"
	"identifEye/entity"
	"identifEye/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RegisterFaceHandler(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": STATUS_ERROR, "message": "Form data required"})
		return
	}

	files, ok := c.Request.MultipartForm.File["faces"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": STATUS_ERROR, "message": "Face image required"})
		return
	}

	dateTime := time.Now().Format("2006-01-02")
	saveDir := filepath.Join(os.Getenv("IMAGE_PATH"), "register", dateTime, c.GetString("username"))

	for i, file := range files {
		err := saveFile(file, saveDir, fmt.Sprintf("%d.jpg", i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to process faces"})
			return
		}
	}

	newUser := entity.User{
		Username: c.GetString("username"),
		Password: c.GetString("password"),
		Name:     c.GetString("name"),
		Email:    c.GetString("email"),
	}
	if result := database.Get().Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to create new user"})
		return
	}

	payloads := jwt.MapClaims{"id": newUser.ID, "face": true}
	token, err := utils.GenerateJWT(payloads)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": STATUS_SUCCESS, "token": token})
}
