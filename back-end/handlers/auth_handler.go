// Обработчик авторизации пользователя
package handlers

import (
	"bytes"
	"dao_vote/back-end/repository"
	"dao_vote/back-end/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// AuthRequest представляет структуру запроса для авторизации
type AuthRequest struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	DeviceName string `json:"device_name"`
}

// AuthResponse представляет структуру ответа для авторизации
type AuthResponse struct {
	Token string `json:"token"`
}

// UserMeResponse представляет структуру ответа от внешнего API /auth/me
type UserMeResponse struct {
	Data User `json:"data"`
}

// User представляет структуру пользователя
type User struct {
	ID            int            `json:"id"`
	Login         string         `json:"login"`
	Email         string         `json:"email"`
	Phone         string         `json:"phone"`
	Nick          string         `json:"nick"`
	Locale        string         `json:"locale"`
	Avatar        string         `json:"avatar"`
	Wallet        string         `json:"wallet"`
	Roles         []Role         `json:"roles"`
	Subscriptions []Subscription `json:"subscriptions"`
}

// Role представляет структуру роли пользователя
type Role struct {
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

// Permission представляет структуру разрешения
type Permission struct {
	Name string `json:"name"`
}

// Subscription представляет структуру подписки пользователя
type Subscription struct {
	ID              int     `json:"id"`
	Tag             string  `json:"tag"`
	PlanID          int     `json:"plan_id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	Currency        string  `json:"currency"`
	TrialPeriod     int     `json:"trial_period"`
	TrialInterval   string  `json:"trial_interval"`
	GracePeriod     int     `json:"grace_period"`
	GraceInterval   string  `json:"grace_interval"`
	InvoicePeriod   int     `json:"invoice_period"`
	InvoiceInterval string  `json:"invoice_interval"`
	Tier            int     `json:"tier"`
	StartsAt        string  `json:"starts_at"`
	EndsAt          string  `json:"ends_at"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// UserLoginHandler обрабатывает запрос авторизации пользователя
func UserLoginHandler(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		var authReq AuthRequest
		if err := c.ShouldBindJSON(&authReq); err != nil {
			logrus.Errorf("Invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return err
		}

		// Подготовка запроса к внешнему API
		jsonData, err := json.Marshal(authReq)
		if err != nil {
			logrus.Errorf("Failed to marshal auth request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
			return err
		}

		// Отправка запроса к внешнему API
		resp, err := http.Post("https://backend.ddapps.io/api/v1/auth/user_login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			logrus.Errorf("Failed to send request to external API: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to authentication service"})
			return err
		}
		defer resp.Body.Close()

		// Чтение и обработка ответа от внешнего API
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Failed to read response body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read authentication response"})
			return err
		}

		var authResp AuthResponse
		if err := json.Unmarshal(body, &authResp); err != nil {
			logrus.Errorf("Failed to unmarshal authentication response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process authentication response"})
			return err
		}

		// Логирование полученного токена
		logrus.Infof("Получен токен: %s", authResp.Token)

		// Возвращение ответа клиенту
		c.JSON(http.StatusOK, authResp)
		return nil
	})
}

// UserMeHandler обрабатывает запрос для получения информации о текущем пользователе
func UserMeHandler(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		token := c.GetHeader("Authorization")

		// Добавление префикса "Bearer" к токену
		if token != "" {
			token = "Bearer " + token
		}

		// Подготовка запроса к внешнему API
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://backend.ddapps.io/api/v1/auth/me?with_user_information=1", nil)
		if err != nil {
			logrus.Errorf("Failed to create request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return err
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// Отправка запроса к внешнему API
		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("Failed to send request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user service"})
			return err
		}
		defer resp.Body.Close()

		// Чтение и обработка ответа от внешнего API
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Failed to read response body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user information"})
			return err
		}

		var userMeResp UserMeResponse
		if err := json.Unmarshal(body, &userMeResp); err != nil {
			logrus.Errorf("Failed to unmarshal user information: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user information"})
			return err
		}

		// Сохранение информации о пользователе в хранилище
		repository.SaveUser(repository.User{
			ID:            userMeResp.Data.ID,
			Login:         userMeResp.Data.Login,
			Email:         userMeResp.Data.Email,
			Phone:         userMeResp.Data.Phone,
			Nick:          userMeResp.Data.Nick,
			Locale:        userMeResp.Data.Locale,
			Avatar:        userMeResp.Data.Avatar,
			Wallet:        userMeResp.Data.Wallet,
			Roles:         extractRoles(userMeResp.Data.Roles),
			Subscriptions: extractSubscriptions(userMeResp.Data.Subscriptions),
		})

		// Логирование полученной информации о пользователе
		logrus.Infof("Получена информация о пользователе: %+v", userMeResp.Data)

		// Возвращение ответа клиенту
		c.JSON(resp.StatusCode, userMeResp)
		return nil
	})
}

// GetUserByIDHandler обрабатывает запрос для получения информации о пользователе по его ID
func GetUserByIDHandler(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		id := c.Param("id")

		// Преобразование ID в число
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return err
		}

		// Получение информации о пользователе из хранилища
		user, err := repository.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return err
		}

		// Возвращение информации о пользователе
		c.JSON(http.StatusOK, user)
		return nil
	})
}

// Вспомогательная функция для извлечения ролей пользователя
func extractRoles(roles []Role) []string {
	var result []string
	for _, role := range roles {
		result = append(result, role.Name)
	}
	return result
}

// Вспомогательная функция для извлечения подписок пользователя
func extractSubscriptions(subscriptions []Subscription) []repository.Subscription {
	var result []repository.Subscription
	for _, sub := range subscriptions {
		result = append(result, repository.Subscription{
			ID:              sub.ID,
			Tag:             sub.Tag,
			PlanID:          sub.PlanID,
			Name:            sub.Name,
			Description:     sub.Description,
			Price:           sub.Price,
			Currency:        sub.Currency,
			TrialPeriod:     sub.TrialPeriod,
			TrialInterval:   sub.TrialInterval,
			GracePeriod:     sub.GracePeriod,
			GraceInterval:   sub.GraceInterval,
			InvoicePeriod:   sub.InvoicePeriod,
			InvoiceInterval: sub.InvoiceInterval,
			Tier:            sub.Tier,
			StartsAt:        sub.StartsAt,
			EndsAt:          sub.EndsAt,
			CreatedAt:       sub.CreatedAt,
			UpdatedAt:       sub.UpdatedAt,
		})
	}
	return result
}

// AuthMiddleware проверяет JWT токен и извлекает информацию о пользователе
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort()
			return
		}

		// Добавление префикса "Bearer" к токену
		if token != "" {
			token = "Bearer " + token
		}

		// Подготовка запроса к внешнему API для проверки токена и получения информации о пользователе
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://backend.ddapps.io/api/v1/auth/me?with_user_information=1", nil)
		if err != nil {
			logrus.Errorf("Failed to create request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// Отправка запроса к внешнему API
		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("Failed to send request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to authentication service"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Чтение и обработка ответа от внешнего API
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Failed to read response body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read authentication response"})
			c.Abort()
			return
		}

		var userMeResp UserMeResponse
		if err := json.Unmarshal(body, &userMeResp); err != nil {
			logrus.Errorf("Failed to unmarshal user information: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user information"})
			c.Abort()
			return
		}

		// Сохранение информации о пользователе в контексте запроса
		c.Set("user", userMeResp.Data)

		// Логирование полученной информации о пользователе
		logrus.Infof("Пользователь авторизован: %+v", userMeResp.Data)

		c.Next()
	}
}

// Вспомогательная функция для проверки роли "admin" (нечувствительная к регистру)
func isAdmin(user User) bool {
	for _, role := range user.Roles {
		if strings.ToLower(role.Name) == "admin" {
			return true
		}
	}
	return false
}
