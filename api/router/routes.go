package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func initializeRoutes(router *gin.Engine) {

	//	handler.InitializeHandler()

	var v1 *gin.RouterGroup = router.Group("api/v1")
	{

		v1.GET("/user", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "get user",
			})
		})
		v1.GET("/users")
		v1.POST("/user")
		v1.DELETE("/user")
		v1.PUT("/user")

	}

}
