package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
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
	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Printf("database connection pool established")

	// Declare an instance of the application struct, containing the config struct and
	app := &application{
		config: cfg,
		logger: logger,
	}

	//Setup server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
