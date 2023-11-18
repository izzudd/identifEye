package routes

import (
	"identifEye/database"
	"identifEye/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Secret(c *gin.Context) {
	userID := c.GetInt("id")

	var user entity.User
	result := database.Get().First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": STATUS_FAILED, "message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": STATUS_SUCCESS,
		"data": gin.H{
			"id":      userID,
			"name":    user.Name,
			"message": "super secret page",
		},
	})
}
