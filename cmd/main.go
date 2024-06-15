package main

import (
	"dao_vote/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	// Middleware для логирования запросов
	r.Use(requestLogger())

	// Маршруты для голосования команды DAO
	r.GET("/dao-team-vote-results", handlers.GetDAOTeamVoteResults)

	// Маршруты для пользовательских голосований
	r.POST("/votes", handlers.CreateVoteHandler)
	r.GET("/votes/:id", handlers.GetVoteHandler)
	r.DELETE("/votes/:id", handlers.DeleteVoteHandler)
	r.POST("/votes/:id/vote", handlers.AddUserVoteHandler)
	r.GET("/votes/:id/votes", handlers.GetUserVotesHandler)

	// Маршруты для авторизации
	r.POST("/auth/login", handlers.UserLoginHandler)
	r.GET("/auth/me", handlers.UserMeHandler)

	// Маршруты для снятия средств
	r.POST("/api/v1/withdraw", handlers.WithdrawHandler)

	// Маршруты для Swagger
	r.StaticFS("/swagger", http.Dir("./swagger"))

	// Получаем порт из переменной окружения, если не указан, используем 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
	if err := r.Run(":" + port); err != nil {
		logrus.Fatalf("Не удалось запустить сервер: %v", err)
	}
}

// requestLogger - это функция middleware, которая логирует детали каждого запроса.
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("Входящий запрос")
		c.Next()
	}
}
