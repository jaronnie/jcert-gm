package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaronnie/jcert-gm/public"
	"github.com/jaronnie/jcert-gm/server/api"
	"github.com/jaronnie/jcert-gm/server/static"
)

// 解决跨域问题
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 必须，指定允许的域名
		c.Header("Access-Control-Allow-Origin", "*")
		// 可选，指定允许的请求方式
		c.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS,PUT,PATCH")
		// 可选，指定自定义 header 参数，多个用 , 隔开
		c.Header("Access-Control-Allow-Headers", "Token")
		// 可选，指定是否允许携带 cookie
		c.Header("Access-Control-Allow-Credentials", "true")
		// 可选，指定时间内减少发送「预检」请求
		c.Header("Access-Control-Max-Age", "60")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func RunServer() {
	e := gin.Default()
	e.Use(Cors())
	// redirect 到 /ui
	e.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/gen")
	})

	gen := e.Group("/gen")
	static.Static(gen, public.Public)

	apiv1 := e.Group("/api")
	api.ApiRouter(apiv1)

	e.Run(":9999")
}
