package main

import (
	"api/Routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	Routes.SetupRouter(router)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
