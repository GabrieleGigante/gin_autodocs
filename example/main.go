package main

import (
	"net/http"

	"github.com/GabrieleGigante/gin_autodocs"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	gin_autodocs.DocumentApi(r, "/api-docs", gin_autodocs.ApiDocumentation{
		Info: gin_autodocs.Info{
			Title: "My test API",
		},
	})
	r.Run(":8080")
}
