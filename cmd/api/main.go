package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	loggerPackage "go-api/internal/logger"
	"go-api/internal/mailer"
	"go-api/internal/models"
	"go-api/internal/utils"
	"os"
	"strings"
)

func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Read the flags into the config struct
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenCons, "db-max-open-cons", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleCons, "db-max-idle-cons", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	smtpPort, _ := utils.GetenvInt("SMTP_PORT")
	flag.IntVar(&cfg.smtp.port, "smtp-port", smtpPort, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")
	cfg.cors.trustedOrigins = strings.Split(os.Getenv("CORS_TRUSTED_ORIGINS"), ",")
	flag.Parse()

	// Initialize a new logger which writes any messages *at or above* the INFO
	// severity level to the standard out stream.
	logger := loggerPackage.New(os.Stdout, loggerPackage.LevelInfo)

	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err, nil)
	}

	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal(err, nil)
		}
	}(db)

	// log a message to say that the connection pool has been successfully established.
	logger.Info("database connection pool established", nil)

	// Declare an instance of the application struct, containing the config struct and
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
		mailer: mailer.New(
			cfg.smtp.host,
			cfg.smtp.port,
			cfg.smtp.username,
			cfg.smtp.password,
			cfg.smtp.sender,
		),
	}

	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.Fatal(err, nil)
	}
}
