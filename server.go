package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rconway/gogithub/routes"
	"github.com/rconway/gogithub/utils"
)

func Server() {
	config := utils.GetConfig()

	engine := gin.Default()
	root := &engine.RouterGroup

	routes.Root(root)
	routes.Info(root)
	routes.Login(root)

	engine.Run(config.ListenAddress) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
