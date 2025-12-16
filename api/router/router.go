package router

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

/*var (
	port = "8080"
	host = "localhost"
)*/

func Initializer() {
	port := os.Getenv("PORT_REST")
	host := os.Getenv("REST_HOST")

	routerHost := fmt.Sprintf("%s:%s", host, port)

	router := gin.Default()

	initializeRoutes(router)
	err := router.Run(routerHost)
	if err != nil {
		log.Fatal("Is not possible run server", err)
	}
	fmt.Printf("Server started on %s", routerHost)
}
