package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func timeMV() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()
		ctx.Next()
		timeElapsed := time.Since(begin)
		log.Printf("%s used time %v", ctx.Request.URL.Path, timeElapsed.Milliseconds())
	}
}

func limitMV() gin.HandlerFunc {
	limitCh := make(chan struct{}, 10)
	return func(ctx *gin.Context) {
		limitCh <- struct{}{}
		log.Printf("current limitCh length: %d", len(limitCh))
		ctx.Next()
		<-limitCh
	}
}

// func main() {
// 	engine := gin.Default()

// 	engine.Use(timeMV())
// 	engine.Use(limitMV())

// 	engine.GET("/boy", func(ctx *gin.Context) {
// 		ctx.String(http.StatusOK, "Hi Boy")
// 	})

// 	engine.Run(":8080")

// }
