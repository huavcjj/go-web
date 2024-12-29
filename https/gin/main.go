package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

func main() {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:   true,
		SSLRedirect: true,
		SSLHost:     "localhost:4000",
	})

	secureFunc := func() gin.HandlerFunc {
		return func(ctx *gin.Context) {
			err := secureMiddleware.Process(ctx.Writer, ctx.Request)
			if err != nil {
				ctx.Abort()
				return
			}
			if status := ctx.Writer.Status(); status > 300 && status < 399 {
				ctx.Abort()
				return
			}
		}
	}()

	engine := gin.Default()
	engine.Use(secureFunc)

	engine.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, TLS!\n")
	})

	engine.RunTLS(":4000", "../cert.pem", "../key.pem")
}
