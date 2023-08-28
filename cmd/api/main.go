package main

import (
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	loggerPackage "go-api/internal/logger"
	"go-api/internal/mailer"
	"go-api/internal/models"
	"go-api/internal/utils"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	// Initialize config struct.
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
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	// If the version flag value is true,
	//then print out the version number and immediately exit.
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	// Initialize logger
	logger := loggerPackage.New(os.Stdout, loggerPackage.LevelInfo)

	// Setup DB connection
	db, err := models.OpenDB(models.DBConfig{
		Dsn:         cfg.DB.Dsn,
		MaxOpenCons: cfg.DB.MaxOpenCons,
		MaxIdleCons: cfg.DB.MaxIdleCons,
		MaxIdleTime: cfg.DB.MaxIdleTime,
	})
	if err != nil {
		logger.Fatal(err, nil)
	}
	// Defer a call to db.Close() so that the connection pool is closed before exit
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal(err, nil)
		}
	}(db)

	// log a message to say that the connection pool has been successfully established.
	logger.Info("database connection pool established", nil)

	// Initialize application
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
		mailer: mailer.New(
			cfg.Smtp.Host,
			cfg.Smtp.Port,
			cfg.Smtp.Username,
			cfg.Smtp.Password,
			cfg.Smtp.Sender,
		),
	}

	// Publish a new "version" variable in the expvar handler containing our application
	expvar.NewString("version").Set(version)
	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))
	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	// Start the server
	err = app.serve()
	if err != nil {
		logger.Fatal(err, nil)
	}
}
