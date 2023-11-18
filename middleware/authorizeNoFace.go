package middleware

import (
	"fmt"
	"identifEye/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthorizeNoFace(c *gin.Context) {
	// Get the token from the request header, query parameter, or wherever it is passed
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": routes.STATUS_ERROR, "message": "Token is missing"})
		c.Abort()
		return
	}

	// Remove the "Bearer " prefix from the token string
	tokenString = tokenString[len("Bearer "):]

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": routes.STATUS_FAILED, "message": "Cannot parse token"})
		c.Abort()
		return
	}

	// Check if the token is valid
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"status": routes.STATUS_FAILED, "message": "Invalid token"})
		c.Abort()
		return
	}

	// Access the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract claims"})
		c.Abort()
		return
	}

	// Pass the extracted claims to the next handler
	c.Set("id", claims["id"])
	c.Set("face", claims["face"])

	c.Next()
}
