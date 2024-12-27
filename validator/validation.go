package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Student struct {
	Name       string    `form:"name" binding:"required"`
	Score      int       `form:"score" binding:"gt=0"`
	Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`
	Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"`
}

var beforeToday validator.Func = func(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now().Truncate(24 * time.Hour)
		return date.Before(today)
	}
	return false
}

func processErr(err error) string {
	if err == nil {
		return ""
	}
	if invalid, ok := err.(*validator.InvalidValidationError); ok {
		return invalid.Error()
	}
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, len(validationErrs))
		for i, err := range validationErrs {
			msgs[i] = fmt.Sprintf("field %s: %s", err.Field(), err.Tag())
		}
		return strings.Join(msgs, ";")
	}
	return "unknown validation error"
}

func main() {
	engine := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("before_today", beforeToday)
	}

	engine.GET("/", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBind(&stu); err != nil {
			msg := processErr(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"student": stu})
	})
	engine.Run(":8080")
}
