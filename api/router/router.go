package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	port = "8080"
	host = "localhost"
)

func Initializer() {
	router := gin.Default()

	initializeRoutes(router)
	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatal("Is not possible run server", err)
	}
	fmt.Println("Server started on localhost:8080")
}
