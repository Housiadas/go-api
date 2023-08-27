package config

import (
	"flag"
	"go-api/internal/utils"
	"os"
	"strings"
)

// Config Define a config struct to hold all the configuration settings for our application.
type Config struct {
	Port int
	Env  string
	DB   struct {
		Dsn         string
		MaxOpenCons int
		MaxIdleCons int
		MaxIdleTime string
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

func New() Config {
	var cfg Config

	// Read the flags
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.DB.Dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenCons, "db-max-open-cons", 50, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleCons, "db-max-idle-cons", 50, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", "20s", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	smtpPort, _ := utils.GetenvInt("SMTP_PORT")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", smtpPort, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")
	cfg.Cors.TrustedOrigins = strings.Split(os.Getenv("CORS_TRUSTED_ORIGINS"), ",")
	flag.Parse()

	return cfg
}
