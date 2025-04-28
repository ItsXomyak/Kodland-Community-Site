package config

import (
	"time"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type SecurityConfig struct {
	MaxLoginAttempts int
	LockoutDuration  time.Duration
	PasswordMinLength int
	RequireTwoFactor bool
}

type Config struct {
	ServerPort    string
	Database      DatabaseConfig
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration
	Security      SecurityConfig
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "12345678")
	viper.SetDefault("DB_NAME", "auth-service")

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("JWT_EXPIRY", "24h")
	viper.SetDefault("REFRESH_EXPIRY", "168h")

	// Настройки безопасности
	viper.SetDefault("MAX_LOGIN_ATTEMPTS", 5)
	viper.SetDefault("LOCKOUT_DURATION", "15m")
	viper.SetDefault("PASSWORD_MIN_LENGTH", 8)
	viper.SetDefault("REQUIRE_TWO_FACTOR", false)

	if err := viper.ReadInConfig(); err != nil {
	}

	cfg := &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
		},
		SMTPHost:      viper.GetString("SMTP_HOST"),
		SMTPPort:      viper.GetInt("SMTP_PORT"),
		SMTPUsername:  viper.GetString("SMTP_USERNAME"),
		SMTPPassword:  viper.GetString("SMTP_PASSWORD"),
		JWTSecret:     viper.GetString("JWT_SECRET"),
		JWTExpiry:     viper.GetDuration("JWT_EXPIRY"),
		RefreshExpiry: viper.GetDuration("REFRESH_EXPIRY"),
		Security: SecurityConfig{
			MaxLoginAttempts: viper.GetInt("MAX_LOGIN_ATTEMPTS"),
			LockoutDuration:  viper.GetDuration("LOCKOUT_DURATION"),
			PasswordMinLength: viper.GetInt("PASSWORD_MIN_LENGTH"),
			RequireTwoFactor: viper.GetBool("REQUIRE_TWO_FACTOR"),
		},
	}
	return cfg, nil
}
