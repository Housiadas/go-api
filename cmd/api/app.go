package main

import (
	"fmt"
	"go-api/internal/logger"
	"go-api/internal/mailer"
	"go-api/internal/models"
	"sync"
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

type application struct {
	config Config
	logger *logger.Logger
	models models.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

// Define an envelope type for the JSON responses
type envelope map[string]interface{}

var (
	buildTime string
	version   string
)

// Background helper callback
func (app *application) background(fn func()) {
	// Increment the WaitGroup counter.
	app.wg.Add(1)
	// Launch a background goroutine.
	go func() {
		// Use defer to decrement the WaitGroup counter before the goroutine returns.
		defer app.wg.Done()
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Errorf("%s", err), nil)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}
