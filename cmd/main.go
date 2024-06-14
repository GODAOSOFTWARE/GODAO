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
	r.POST("/votes", handlers.CreateVoteHandler)            //Cоздает голосование пользователей
	r.GET("/votes/:id", handlers.GetVoteHandler)            //Получает информацию о голосвании по ID
	r.DELETE("/votes/:id", handlers.DeleteVoteHandler)      //Удаляет голосование по ID
	r.POST("/votes/:id/vote", handlers.AddUserVoteHandler)  // Создает голос пользователя в голосовании
	r.GET("/votes/:id/votes", handlers.GetUserVotesHandler) // Получает голоса пользователей по ID голосования

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
