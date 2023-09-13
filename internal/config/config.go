package config

import (
	"go-api/internal/utils"
	"os"
	"strings"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		Dsn         string
		MaxOpenCons int
		MaxIdleCons int
		MaxIdleTime string
	}
	Redis struct {
		Address  string
		Password string
	}
	// Add a new limiter struct containing fields for the requests-per-second and burst
	// values, and a boolean field which we can use to enable/disable rate limiting
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	Cors struct {
		TrustedOrigins []string
	}
}

func NewConfig() Config {
	var cfg Config
	// APP
	cfg.Env = os.Getenv("APP_ENV")
	cfg.Port, _ = utils.GetenvInt("APP_PORT")
	// Database
	cfg.DB.Dsn = os.Getenv("DB_DSN")
	cfg.DB.MaxOpenCons, _ = utils.GetenvInt("DB_MAX_OPEN_CONNECTIONS")
	cfg.DB.MaxIdleCons, _ = utils.GetenvInt("DB_MAX_IDLE_CONNECTIONS")
	cfg.DB.MaxIdleTime = os.Getenv("DB_MAX_CONNECTION_IDLE_TIME")
	// Redis
	cfg.Redis.Address = os.Getenv("REDIS_ADDR")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	// Rate Limiter
	cfg.Limiter.Rps = 2 // requests per second
	cfg.Limiter.Burst = 4
	cfg.Limiter.Enabled, _ = utils.GetenvBool("LIMITER_ENABLED")
	// SMTP
	cfg.Smtp.Host = os.Getenv("SMTP_HOST")
	cfg.Smtp.Port, _ = utils.GetenvInt("SMTP_PORT")
	cfg.Smtp.Username = os.Getenv("SMTP_USERNAME")
	cfg.Smtp.Password = os.Getenv("SMTP_PASSWORD")
	cfg.Smtp.Sender = os.Getenv("SMTP_SENDER")
	// CORS
	cfg.Cors.TrustedOrigins = strings.Split(os.Getenv("CORS_TRUSTED_ORIGINS"), ",")
	return cfg
}
