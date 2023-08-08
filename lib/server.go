package lib

import "github.com/gin-gonic/gin"

var GinServer *gin.Engine
var SysConfig *Config

func RunServer() {
	GinServer = gin.New()

	SysConfig = loadConfigs()
	registerRoutes()

	GinServer.Use(func(ctx *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				ctx.AbortWithStatusJSON(400, gin.H{"error": e})
			}
		}()
		ctx.Next()
	})

	GinServer.POST("/route", func(ctx *gin.Context) {
		rc := &RouteConfig{}
		if err := ctx.ShouldBindJSON(rc); err != nil {
			panic(err)
		}
		rc.Register(true)
		ctx.JSON(200, gin.H{"msg": "success"})
	})

	GinServer.GET("/config", func(ctx *gin.Context) {
		ctx.JSON(200, SysConfig)
	})

	GinServer.Run(":8080")
}
