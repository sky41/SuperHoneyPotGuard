package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GinMode         string
	DBHost          string
	DBPort          string
	DBName          string
	DBUser          string
	DBPassword      string
	JWTSecret       string
	JWTExpiresIn    string
	BCryptCost      int
	RateLimitWindow int
	RateLimitMax    int
	LogLevel        string
	LogFilePath     string
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	RedisDB         int
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = &Config{
		Port:            getEnv("PORT", "3000"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBName:          getEnv("DB_NAME", "superhoneypotguard"),
		DBUser:          getEnv("DB_USER", "root"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		JWTSecret:       getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
		JWTExpiresIn:    getEnv("JWT_EXPIRES_IN", "24h"),
		BCryptCost:      getEnvAsInt("BCRYPT_COST", 10),
		RateLimitWindow: getEnvAsInt("RATE_LIMIT_WINDOW_MS", 900000),
		RateLimitMax:    getEnvAsInt("RATE_LIMIT_MAX_REQUESTS", 100),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		LogFilePath:     getEnv("LOG_FILE_PATH", "logs/"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnvAsInt("REDIS_DB", 0),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal := parseInt(value); intVal != 0 {
			return intVal
		}
	}
	return defaultValue
}

func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return result
}
