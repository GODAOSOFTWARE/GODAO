package handlers

import (
	"bytes"
	"dao_vote/internal/repository"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
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
	ID     int    `json:"id"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Nick   string `json:"nick"`
	Locale string `json:"locale"`
	Avatar string `json:"avatar"`
	Wallet string `json:"wallet"`
	Roles  []struct {
		Name        string `json:"name"`
		Permissions []struct {
			Name string `json:"name"`
		} `json:"permissions"`
	} `json:"roles"`
	Subscriptions []struct {
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
	} `json:"subscriptions"`
}

// UserLoginHandler обрабатывает запрос авторизации пользователя
func UserLoginHandler(c *gin.Context) {
	var authReq AuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Подготовка запроса к внешнему API
	jsonData, err := json.Marshal(authReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправка запроса к внешнему API
	resp, err := http.Post("https://backend.ddapps.io/api/v1/auth/user_login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Чтение и обработка ответа от внешнего API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращение ответа клиенту
	c.JSON(http.StatusOK, authResp)
}

// UserMeHandler обрабатывает запрос для получения информации о текущем пользователе
func UserMeHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")

	// Добавление префикса "Bearer" к токену
	if token != "" {
		token = "Bearer " + token
	}

	// Подготовка запроса к внешнему API
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://backend.ddapps.io/api/v1/auth/me?with_user_information=1", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Отправка запроса к внешнему API
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Чтение и обработка ответа от внешнего API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userMeResp UserMeResponse
	if err := json.Unmarshal(body, &userMeResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	// Возвращение ответа клиенту
	c.JSON(resp.StatusCode, userMeResp)
}

// GetUserByIDHandler обрабатывает запрос для получения информации о пользователе по его ID
func GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")

	// Преобразование ID в число
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Получение информации о пользователе из хранилища
	user, err := repository.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Возвращение информации о пользователе
	c.JSON(http.StatusOK, user)
}

// Вспомогательная функция для извлечения ролей пользователя
func extractRoles(roles []struct {
	Name        string `json:"name"`
	Permissions []struct {
		Name string `json:"name"`
	} `json:"permissions"`
}) []string {
	var result []string
	for _, role := range roles {
		result = append(result, role.Name)
	}
	return result
}

// Вспомогательная функция для извлечения подписок пользователя
func extractSubscriptions(subscriptions []struct {
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
}) []repository.Subscription {
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
