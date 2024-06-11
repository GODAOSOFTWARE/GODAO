package utils

import (
	"github.com/gin-gonic/gin"
)

func JSONResponse(c *gin.Context, code int, message interface{}) {
	c.JSON(code, message)
}
