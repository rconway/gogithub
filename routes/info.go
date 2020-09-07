package routes

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/rconway/gogithub/middleware"
)

const infoPage = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Document</title>
  </head>
  <body>
		<div>
			Hello {{.Name}}. How are you?
    </div>
  </body>
</html>
`

// Info zzz
func Info(r *gin.RouterGroup) *gin.RouterGroup {
	info := r.Group("info")

	info.Use(middleware.EnsureUser())

	info.GET("", func(c *gin.Context) {
		tpl := template.Must(template.New("info").Parse(infoPage))
		userEmail, _ := c.Get("user-email")
		data := struct {
			Name string
		}{}
		data.Name = userEmail.(string)
		tpl.ExecuteTemplate(c.Writer, "info", data)
	})

	return info
}
