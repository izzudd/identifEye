package routes

import (
	"identifEye/database"
	"identifEye/entity"
	"identifEye/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var body loginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": STATUS_ERROR, "message": "Required field: username, password"})
		return
	}

	user := entity.User{}
	if err := database.Get().Where("username = ?", body.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"status": STATUS_ERROR, "message": "Invalid credentials"})
		return
	}

	if !utils.PasswordMatch(body.Password, user.Password) {
		c.JSON(http.StatusConflict, gin.H{"status": STATUS_ERROR, "message": "Invalid credentials"})
		return
	}

	bodyMap := jwt.MapClaims{"id": user.ID, "face": false}
	token, err := utils.GenerateJWT(bodyMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": STATUS_SUCCESS, "token": token})
}
