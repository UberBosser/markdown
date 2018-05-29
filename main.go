package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"html/template"
	"io/ioutil"
	"os"
)

type MarkdownFile struct {
	Name string        `json:"name"`
	Html template.HTML `json:"html"`
}

func render(c *gin.Context) {
	c.HTML(200, "main.tmpl", nil)
}

func data(c *gin.Context) {
	md, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		c.JSON(200, template.HTML("<h2>Error</h2><p>"+err.Error()+"</p>"))
		return
	}
	markdownFile := MarkdownFile{os.Args[1], template.HTML(markdown.ToHTML(md, nil, nil))}
	c.JSON(200, markdownFile)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough args, markdown file should be an argument")
		return
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.LoadHTMLGlob("./templates/main.tmpl")

	router.GET("/", render)
	router.GET("/data", data)

	router.NoRoute(render)

	fmt.Println("Serving port :8080")
	router.Run()
}
