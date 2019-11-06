package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"html/template"
	"io/ioutil"
	"os"
)

var mainTemplate = template.Must(template.New("main").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>MdServer</title>
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.10.0/github-markdown.css" />
</head>
<body class="markdown-body">
	{{ . }}
	<script
		src="https://code.jquery.com/jquery-3.3.1.min.js"
		integrity="sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8="
		crossorigin="anonymous">
	</script>
	<script>function update(){$.ajax({type:"GET",url:"/data",success:function(t){$("body").html(t.html);if(document.title!=t.name){document.title=t.name;}}})}$(document).ready(function(){update(),setInterval(update,1e3)});</script>
</body>
</html>`))

type MarkdownFile struct {
	Name string        `json:"name"`
	Html template.HTML `json:"html"`
}

func render(c *gin.Context) {
	c.HTML(200, "main", nil)
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
	router.SetHTMLTemplate(mainTemplate)

	router.GET("/", render)
	router.GET("/data", data)

	router.NoRoute(render)

	fmt.Println("Serving " + os.Args[1] + " to port :8080")
	router.Run(":8080")
}
