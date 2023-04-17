package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Recovery(c *gin.Context) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("Was catch panic in panic middleware")
			c.JSON(http.StatusInternalServerError, ErrorResponse(errors.New("internal Server error")))
		}
	}()
	c.Next()
}
