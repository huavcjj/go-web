package main

import (
	"log"
	"net/http"
	"strconv"

	_ "go-web/swagger/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//	@Summary	Get user by ID
//	@Produce	json
//	@Param		id	path		int		true	"User ID"
//	@Success	200	{object}	User	"OK"
//	@Failure	400	{object}	string	"Invalid ID"
//	@Failure	500	{object}	string	"Internal Server Error"
//	@Router		/user/{id} [get]
func GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if id, err := strconv.Atoi(idStr); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid ID")
	} else {
		ctx.JSON(http.StatusOK, User{
			ID:   id,
			Name: "John Doe",
			Age:  25,
		})
	}
}

//	@Summary	Update user
//	@Produce	json
//	@Param		user	body		User	true	"User object"
//	@Success	200		{object}	User	"OK"
//	@Failure	400		{object}	string	"Invalid request"
//	@Failure	500		{object}	string	"Internal Server Error"
//	@Router		/user [post]
func UpdateUser(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.String(http.StatusBadRequest, "Invalid request")
	} else {
		ctx.JSON(http.StatusOK, user)
	}
}

func main() {
	engine := gin.Default()

	engine.GET("/swagger/*all", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.GET("/user/:id", GetUser)
	engine.POST("/user", UpdateUser)

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}

//swag fmt main.go
//swag init -g main.go
