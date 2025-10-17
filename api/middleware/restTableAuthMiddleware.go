package middleware

import (
	"net/http"
	"strconv"

	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RestTableAuthMiddleware(db *gorm.DB, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		tableIDStr := c.Param(paramName)
		tableID, _ := strconv.Atoi(tableIDStr)

		var tableUserModel models.TableUser

		if err := db.Where("user_id = ? AND table_id = ?", userID, tableID).First(&tableUserModel).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "permission denied"})
			return
		}
		c.Next()
	}
}
