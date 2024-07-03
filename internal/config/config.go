package config

import "os"

var (
	AuthServiceURL = getEnv("AUTH_SERVICE_URL", "https://backend.ddapps.io/api/v1/auth")
	WithdrawURL    = getEnv("WITHDRAW_URL", "https://backend.ddapps.io/api/v1/withdraw")
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
