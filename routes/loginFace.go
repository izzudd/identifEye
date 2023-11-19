package routes

import (
	"bytes"
	"fmt"
	"identifEye/database"
	"identifEye/entity"
	"identifEye/utils"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	saveDir := filepath.Join(os.Getenv("IMAGE_PATH"), "login", fmt.Sprint(c.GetInt("id")), dateTime)

	for i, file := range files {
		err := saveFile(file, saveDir, fmt.Sprintf("%d.jpg", i))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to process faces"})
			return
		}
	}

	similarity, err := getSimilarityScore(c.GetInt("id"), saveDir)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": STATUS_FAILED, "message": "Face not detected"})
		return
	}

	treshold, _ := strconv.ParseFloat(os.Getenv("SIMILARITY_TRESHOLD"), 32)
	if similarity < float32(treshold) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": STATUS_FAILED, "message": "Face does not match", "data": gin.H{
			"similarity": similarity,
		}})
		return
	}

	payloads := jwt.MapClaims{"id": c.GetInt("id"), "face": true}
	token, err := utils.GenerateJWT(payloads)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": STATUS_SUCCESS, "token": token, "data": gin.H{
		"similarity": similarity,
	}})
}

func getSimilarityScore(userId int, imagePath string) (float32, error) {
	var user entity.User
	database.Get().First(&user, userId)

	scriptPath := filepath.Join(".", "model", "detect.py")
	imageAbsolutePath, _ := filepath.Abs(imagePath)

	cmd := exec.Command("python", scriptPath, user.Key, imageAbsolutePath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return -1, fmt.Errorf("%s, %s", err.Error(), stdout.String())
	}

	stringNum := strings.TrimRight(strings.Split(stdout.String(), ": ")[1], "\n")
	similarity, err := strconv.ParseFloat(stringNum, 32)
	return float32(similarity), err
}
