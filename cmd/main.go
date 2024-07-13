package main

import (
	"dao_vote/back-end/handlers"
	"dao_vote/back-end/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	// Инициализация базы данных
	if err := repository.InitDB("./votes.db"); err != nil {
		logrus.Fatalf("Не удалось инициализировать базу данных: %v", err)
	}

	// Применение миграций
	applyMigrations()

	r := setupRouter() // Настраиваем маршруты

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

func applyMigrations() {
	m, err := migrate.New(
		"file://migrations",
		"sqlite3://./votes.db")
	if err != nil {
		logrus.Fatalf("Ошибка инициализации миграции: %v", err)
	}

	// Проверка и установка версии базы данных
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		logrus.Fatalf("Ошибка получения версии базы данных: %v", err)
	}

	if dirty {
		if err := m.Force(int(version)); err != nil {
			logrus.Fatalf("Ошибка принудительной установки версии базы данных: %v", err)
		}
		logrus.Infof("Принудительно установлена версия базы данных: %v", version)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Ошибка применения миграций: %v", err)
	}

	logrus.Info("Миграции успешно применены")
}

// setupRouter - это функция, которая настраивает маршруты и возвращает экземпляр gin.Engine.
func setupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware для логирования запросов
	r.Use(requestLogger())

	// Middleware для CORS
	r.Use(corsMiddleware())

	// Маршруты для Swagger (должны быть доступны без авторизации)
	r.StaticFS("/swagger", http.Dir("./swagger"))

	// Группа маршрутов, которые требуют авторизации
	authRoutes := r.Group("/")
	authRoutes.Use(handlers.AuthMiddleware())
	{
		// Маршруты для голосования команды DAO
		authRoutes.GET("/get-voting-results-by-wallet", handlers.GetVotingResultsByWallet)

		// Маршруты для пользовательских голосований
		authRoutes.POST("/votes", handlers.CreateVoteHandler)
		authRoutes.GET("/votes/:id", handlers.GetVoteHandler)
		authRoutes.DELETE("/votes/:id", handlers.DeleteVoteHandler)
		authRoutes.POST("/votes/:id/vote", handlers.AddUserVoteHandler)
		authRoutes.GET("/votes/:id/votes", handlers.GetUserVotesHandler)

		// Маршруты для снятия средств
		authRoutes.POST("/api/v1/withdraw", handlers.WithdrawHandler)

		// Маршруты для получения названий таблиц и элементов в таблице
		authRoutes.GET("/tables", handlers.GetTableNamesHandler)
		authRoutes.GET("/tables/:table_name/elements", handlers.GetTableElementsHandler)

		// Маршруты для администрирования кошельков
		authRoutes.POST("/wallets", handlers.AddWalletHandler)
		authRoutes.DELETE("/wallets/:wallet_address", handlers.DeleteWalletHandler)
	}

	// Маршруты для авторизации (не требуют авторизации)
	r.POST("/auth/login", handlers.UserLoginHandler)
	r.GET("/auth/me", handlers.UserMeHandler)

	return r
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
