// Package utils Утилита для оптимизации обработчиков
package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// HandleRequest оборачивает обработку запроса, добавляя логирование ошибок
func HandleRequest(c *gin.Context, handlerFunc func(c *gin.Context) error) {
	if err := handlerFunc(c); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
