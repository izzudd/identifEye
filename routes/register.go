package routes

import (
	"identifEye/database"
	"identifEye/entity"
	"identifEye/utils"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type registerBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

func RegisterHandler(c *gin.Context) {
	var body registerBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": STATUS_ERROR, "message": "Required field: username, password, email, name"})
		return
	}

	existingUser := entity.User{}
	if err := database.Get().Where("username = ?", body.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": STATUS_ERROR, "message": "User with the same username already exists"})
		return
	}

	bodyMap := structToMap(body)
	hasedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to process request"})
		return
	}
	bodyMap["password"] = hasedPassword

	token, err := utils.GenerateJWT(bodyMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": STATUS_ERROR, "message": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": STATUS_SUCCESS, "token": token})
}

func structToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Tag.Get("json")
		fieldValue := v.Field(i).Interface()
		result[fieldName] = fieldValue
	}

	return result
}
