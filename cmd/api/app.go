package main

import (
	"fmt"
	"go-api/internal/cache"
	"go-api/internal/config"
	"go-api/internal/logger"
	"go-api/internal/mailer"
	"go-api/internal/models"
	"sync"
)

type application struct {
	config config.Config
	logger *logger.Logger
	models models.Models
	cache  cache.Cache
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

// Define an envelope type for the JSON responses
type envelope map[string]interface{}

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
