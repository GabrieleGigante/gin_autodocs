package gin_autodocs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

var docs map[string]Operation = make(map[string]Operation)

// func trace(depth int) (string, int, string) {
// 	pc, file, line, ok := runtime.Caller(depth)
// 	if !ok {
// 		return "?", 0, "?"
// 	}

// 	fn := runtime.FuncForPC(pc)
// 	if fn == nil {
// 		return file, line, "?"
// 	}

// 	return file, line, fn.Name()
// }

func DocumentEndpoint(endpoint string, doc Operation) {
	// file, line, fnName := trace(2)
	// fmt.Println(fmt.Sprint(file, " | ", line, " | ", fnName))
	docs[endpoint] = doc
}

func DocumentApi(e *gin.Engine, path string, apiDocs DocumentationOptions) {
	mode := os.Getenv("GIN_MODE")
	if mode != "debug" && mode != "" {
		return
	}
	routes := e.Routes()
	for _, route := range routes {
		if apiDocs.Paths == nil {
			apiDocs.Paths = make(map[string]map[string]Operation)
		}
		if apiDocs.Paths[route.Path] == nil {
			apiDocs.Paths[route.Path] = map[string]Operation{}
		}
		temp := docs[strings.ToUpper(route.Method)+"-"+route.Path]
		temp.Tags = append(temp.Tags, GetAutoTag(route.Path))
		fmt.Println(temp)
		apiDocs.Paths[route.Path][strings.ToLower(route.Method)] = temp
	}

	e.GET(path, func(c *gin.Context) {
		c.Header("Content-type", "text/html")
		c.File("docs/docs.html")
	})
	e.GET("/favicon.ico", func(c *gin.Context) {
		c.File("favicon.ico")
	})

	res, err := apiDocs.toJson()
	fmt.Println(string(res))
	if err != nil {
		panic(err)
	}
	// generating docs
	os.Mkdir("docs", os.ModePerm)
	err = os.WriteFile("docs/docs.json", res, 0644)
	if err != nil {
		panic(err.Error())
	}
	out, err := exec.Command("redoc-cli", "build", "docs/docs.json", "-o", "docs/docs.html").Output()
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(string(out))
}

func GetAutoTag(path string) string {
	list := strings.Split(path, "/")
	if len(list) > 1 {
		return list[1]
	}
	return ""
}

func (d DocumentationOptions) toJson() ([]byte, error) {
	// paths := make(map[string]interface{})
	// for _, path := range d.Paths {
	// }
	if d.Openapi == "" {
		d.Openapi = "3.0.0"
	}
	// d.Paths = make(map[string]map[string]Operation)

	res, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return make([]byte, 0), err
	}
	return res, nil
}

type DocumentationOptions struct {
	Openapi      string                          `json:"openapi,omitempty"`
	Info         Info                            `json:"info,omitempty"`
	ExternalDocs *ExternalDocs                   `json:"externalDocs,omitempty"`
	Servers      []Servers                       `json:"servers,omitempty"`
	Tags         []Tags                          `json:"tags,omitempty"`
	Paths        map[string]map[string]Operation `json:"paths,omitempty"`
}
type Info struct {
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version,omitempty"`
}
type Contact struct {
	Email string `json:"email,omitempty"`
}
type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}
type Servers struct {
	URL string `json:"url,omitempty"`
}
type Tags struct {
	Name         string       `json:"name,omitempty"`
	Description  string       `json:"description,omitempty"`
	ExternalDocs ExternalDocs `json:"externalDocs,omitempty"`
}
type Operation struct {
	Tags        []string    `json:"tags,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	OperationID string      `json:"operationId,omitempty"`
	RequestBody RequestBody `json:"requestBody,omitempty"`
	// Responses   Responses   `json:"responses,omitempty"`
}
type RequestBody struct {
	Description string         `json:"description,omitempty"`
	Content     map[string]any `json:"content,omitempty"`
	Required    bool           `json:"required,omitempty"`
}
