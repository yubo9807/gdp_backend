package main

import (
	"gdp/service/configs"
	"gdp/service/controller"
	"gdp/service/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	server := app.Group(configs.Service.Prefix)
	server.Use(middleware.Log)
	server.Use(middleware.ContextMiddleware)

	server.GET("/file/catalog", controller.FileCatalog)
	server.GET("/file/info", controller.FileInfo)
	server.GET("/file/content", controller.FileContent)

	server.POST("/publish", controller.Verification, controller.Publish)
	server.DELETE("/unpublish", controller.Verification, controller.Unpublish)
	server.GET("/install", controller.Install)

	server.POST("/user/register", controller.Register)
	server.PUT("/user/activated", controller.Verification, controller.VerificationRole(0), controller.Activated)
	server.POST("/user/login", controller.Login)
	server.GET("/user/info", controller.Verification, controller.UserInfo)
	server.GET("/user/list", controller.UserList)

	port := ":" + strconv.Itoa(configs.Service.Port)
	app.Run(port)
}
