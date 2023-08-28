package main

import (
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"go-api/internal/config"
	loggerPackage "go-api/internal/logger"
	"go-api/internal/mailer"
	"go-api/internal/models"
	"os"
	"runtime"
	"time"
)

var (
	version   string
	buildTime string
)

func main() {
	// If the version flag value is true,
	// then print out the version number and build time and immediately exit.
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()
	if *displayVersion {
		// Linker flags - check makefile
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	// Initialize config
	cfg := config.NewConfig()

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
