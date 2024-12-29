package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

var (
	userInfos sync.Map
)

const (
	authCookie = "auth"
)

func genSessionID(ctx *gin.Context) string {
	return base64.StdEncoding.EncodeToString([]byte(ctx.Request.RemoteAddr))
}

type User struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Vip  bool   `json:"vip"`
}

func login(ctx *gin.Context) {
	sessionID := genSessionID(ctx)
	user := User{
		Name: "Golang",
		Role: "admin",
		Vip:  true,
	}
	userInfo, _ := sonic.Marshal(user)
	userInfos.Store(sessionID, userInfo)

	ctx.SetCookie(
		authCookie,
		sessionID,
		3000,
		"/",
		"localhost",
		false,
		true,
	)
	ctx.String(http.StatusOK, "login success")
}

func myProfile(ctx *gin.Context) {
	ctx.String(http.StatusOK, fmt.Sprintf("my profile: %v", ctx.GetString("name")))
}

func postVideo(ctx *gin.Context) {
	ctx.String(http.StatusOK, "post video success")
}

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(authCookie)
		if err != nil {
			fmt.Printf("failed to get cookie: %v\n", err)
			ctx.Redirect(http.StatusMovedPermanently, "http://localhost:5678/go_to_login")
			ctx.Abort()
			return
		}

		userInfo, ok := userInfos.Load(cookie)
		if !ok {
			ctx.String(http.StatusForbidden, "invalid session")
			ctx.Abort()
			return
		}

		var user User
		if err := sonic.Unmarshal(userInfo.([]byte), &user); err != nil {
			ctx.String(http.StatusInternalServerError, "failed to unmarshal user info")
			ctx.Abort()
			return
		}

		ctx.Set("name", user.Name)
		ctx.Set("role", user.Role)
		ctx.Set("vip", user.Vip)
		ctx.Next()
	}
}

func main() {
	engine := gin.Default()

	engine.LoadHTMLFiles("../login.html")
	engine.GET("/go_to_login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})

	engine.GET("/login", login)
	engine.GET("/profile", authMiddleware(), myProfile)
	engine.GET("/video", authMiddleware(), postVideo)
	engine.Run(":5678")
}
