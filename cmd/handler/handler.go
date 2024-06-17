package handler

import (
	"dao_vote/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// setupRouter - это функция, которая настраивает маршруты и возвращает экземпляр gin.Engine.
func setupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware для логирования запросов
	r.Use(requestLogger())

	// Middleware для CORS
	r.Use(corsMiddleware())

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

	return r
}

// Handler - экспортированная функция, которую Vercel будет использовать.
func Handler(w http.ResponseWriter, r *http.Request) {
	router := setupRouter()
	router.ServeHTTP(w, r)
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

// corsMiddleware - это функция middleware, которая добавляет необходимые заголовки для разрешения CORS.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
