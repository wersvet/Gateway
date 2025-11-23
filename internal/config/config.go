package config

import "os"

// Config holds application configuration values.
type Config struct {
	Port           string
	AuthServiceURL string
	UserServiceURL string
	ChatServiceURL string
	JWTSecret      string
}

// Load reads configuration from environment variables with sane defaults.
func Load() Config {
	return Config{
		Port:           getEnv("PORT", "9000"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL: getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		ChatServiceURL: getEnv("CHAT_SERVICE_URL", "http://localhost:8083"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
