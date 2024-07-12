// Package models Cтруктура для запросов на перемещение токенов
package models

// WithdrawRequest представляет структуру запроса для снятия средств
type WithdrawRequest struct {
	Amount  float64 `json:"amount" validate:"required"`
	Address string  `json:"address" validate:"required"`
}
