package main

import (
	"dao_vote/internal/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	// Подключаем маршруты для голосования команды DAO
	r.GET("/dao-team-vote-results", handlers.GetDAOTeamVoteResults)

	// Подключаем маршруты для пользовательских голосований
	r.POST("/votes", handlers.CreateVoteHandler)
	r.GET("/votes/:id", handlers.GetVoteHandler)
	r.DELETE("/votes/:id", handlers.DeleteVoteHandler)
	r.POST("/votes/:id/vote", handlers.AddUserVoteHandler)
	r.GET("/votes/:id/votes", handlers.GetUserVotesHandler)

	// Подключаем маршруты для авторизации
	r.POST("/auth/login", handlers.UserLoginHandler) // Добавлено

	// Подключаем маршруты для Swagger
	r.StaticFS("/swagger", http.Dir("./swagger"))

	// Получаем порт из переменной окружения, если не указан, используем 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Запускаем сервер
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
