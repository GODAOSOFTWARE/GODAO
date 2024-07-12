// Package repository Хранилище пользователей
package repository

import (
	"errors"
	"sync"
)

type User struct {
	ID            int
	Login         string
	Email         string
	Phone         string
	Nick          string
	Locale        string
	Avatar        string
	Wallet        string
	Roles         []string
	Permissions   []string
	Subscriptions []Subscription
}

type Subscription struct {
	ID              int
	Tag             string
	PlanID          int
	Name            string
	Description     string
	Price           float64
	Currency        string
	TrialPeriod     int
	TrialInterval   string
	GracePeriod     int
	GraceInterval   string
	InvoicePeriod   int
	InvoiceInterval string
	Tier            int
	StartsAt        string
	EndsAt          string
	CreatedAt       string
	UpdatedAt       string
}

var (
	userStore = make(map[int]User)
	userMutex = sync.Mutex{}
)

// SaveUser сохраняет информацию о пользователе в хранилище
func SaveUser(user User) {
	userMutex.Lock()
	defer userMutex.Unlock()
	userStore[user.ID] = user
}

// GetUserByID возвращает информацию о пользователе по его ID
func GetUserByID(id int) (User, error) {
	userMutex.Lock()
	defer userMutex.Unlock()
	user, exists := userStore[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}
