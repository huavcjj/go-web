package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var JWTSecret = os.Getenv("JWT_SECRET")

var DefaultHeader = JwtHeader{
	Alg: "HS256",
	Typ: "JWT",
}

type JwtHeader struct {
	Alg string `json:"alg"` // 使用するアルゴリズム
	Typ string `json:"typ"` // JWTのタイプ
}

type JwtPayload struct {
	ID          string                 `json:"jti"`          // JWT ID: トークンを一意に識別するためのID
	Issue       string                 `json:"iss"`          // Issuer: トークンを発行したエンティティ（例: 発行者）
	Audience    string                 `json:"aud"`          // Audience: トークンの受信者（例: APIクライアント）
	Subject     string                 `json:"sub"`          // Subject: トークンの主題（例: ユーザーID）
	IssueAt     int64                  `json:"iat"`          // Issued At: トークンが発行された時刻（Unixタイムスタンプ）
	NotBefore   int64                  `json:"nbf"`          // Not Before: トークンが使用可能になる時刻（Unixタイムスタンプ）
	Expiration  int64                  `json:"exp"`          // Expiration Time: トークンの有効期限（Unixタイムスタンプ）
	UserDefined map[string]interface{} `json:"ud,omitempty"` // 任意のユーザー定義フィールド（必要に応じて追加）
}

func GenJWT(header JwtHeader, payload JwtPayload) (string, error) {
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	h := hmac.New(sha256.New, []byte(JWTSecret))
	h.Write([]byte(headerEncoded + "." + payloadEncoded))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return headerEncoded + "." + payloadEncoded + "." + signature, nil
}

func VerifyJWT(token string) (*JwtHeader, *JwtPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, nil, fmt.Errorf("invalid token")
	}

	h := hmac.New(sha256.New, []byte(JWTSecret))
	h.Write([]byte(parts[0] + "." + parts[1]))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if signature != parts[2] {
		return nil, nil, fmt.Errorf("invalid signature")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, err
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, err
	}

	var header JwtHeader
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, nil, err
	}

	var payload JwtPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, nil, err
	}

	return &header, &payload, nil
}

func Login(ctx *gin.Context) {
	header := DefaultHeader
	payload := JwtPayload{
		Issue:      "golang",
		IssueAt:    time.Now().Unix(),
		Expiration: time.Now().Add(time.Hour * 24).Unix(),
		UserDefined: map[string]interface{}{
			"name": "Golang",
			"role": "admin",
			"vip":  true,
		},
	}

	token, err := GenJWT(header, payload)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Header("token", token)
	ctx.String(http.StatusOK, "login success")
}

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		_, payload, err := VerifyJWT(token)
		if err != nil {
			ctx.String(http.StatusForbidden, err.Error())
			ctx.Abort()
			return
		}

		for k, v := range payload.UserDefined {
			ctx.Set(k, v)
		}
		ctx.Next()
	}
}

func main() {
	engine := gin.Default()
	engine.GET("/login", Login)
	engine.GET("/profile", JwtAuthMiddleware(), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, fmt.Sprintf("my profile: %v", ctx.GetString("name")))
	})
	engine.Run(":5555")
}
