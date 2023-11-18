package middleware

import (
	"identifEye/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeFace(c *gin.Context) {
	faceAuth := c.GetBool("face")
	if !faceAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"status": routes.STATUS_FAILED, "message": "Face unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
