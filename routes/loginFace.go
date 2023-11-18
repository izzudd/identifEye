package routes

import (
	"fmt"
	"identifEye/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func LoginFaceHandler(c *gin.Context) {
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

	dateTime := time.Now().Format("2006-01-02_15-04-05")
	saveDir := filepath.Join(fmt.Sprint(1), os.Getenv("IMAGE_PATH"), "login", dateTime)

	for i, file := range files {
		err := saveFile(file, saveDir, fmt.Sprintf("%d.jpg", i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to process faces"})
			return
		}
	}

	payloads := jwt.MapClaims{"id": c.GetString("username"), "face": true}
	token, err := utils.GenerateJWT(payloads)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": STATUS_SUCCESS, "token": token})
}
