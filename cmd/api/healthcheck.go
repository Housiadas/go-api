package main

import (
	"fmt"
	"net/http"
)

// Declare a handler which writes a plain-text response
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)
	// Set the "Content-Type: application/json" header on the response.
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON as the HTTP response body.
	w.Write([]byte(js))
}
