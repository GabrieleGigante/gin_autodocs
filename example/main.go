package main

import (
	"net/http"

	"github.com/GabrieleGigante/gin_autodocs"
	"github.com/gin-gonic/gin"
)

type MyModel struct {
	A string `json:"a"`
	B string `json:"b"`
}

func main() {
	r := gin.Default()
	gin_autodocs.DocumentEndpoint("GET-/hello/world", gin_autodocs.Operation{
		Summary:     "test",
		Description: "my description",
		RequestBody: gin_autodocs.RequestBody{
			Required:    true,
			Description: "Descrizione",
			Content: gin.H{
				"application/json": gin.H{
					"examples": gin.H{
						"myexample": &MyModel{
							A: "A string",
							B: "B string",
						},
					},
				},
			},
		},
	})
	r.GET("/hello/world", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	gin_autodocs.DocumentApi(r, "/api-docs", gin_autodocs.DocumentationOptions{
		Info: gin_autodocs.Info{
			Title:       "My test API",
			Description: "API to test the autodocs",
		},
	})
	// fmt.Println(gin_autodocs.DocumentEndpoint)

	r.Run(":8081")
}
