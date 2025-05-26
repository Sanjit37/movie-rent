package main

import "github.com/gin-gonic/gin"

func main() {
	route := gin.Default()

	route.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"Greetings": "Hello world"})
	})

	route.Run(":8080")
}
