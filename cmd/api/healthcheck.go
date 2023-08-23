package main

import (
	"go-api/internal/serializer"
	"net/http"
)

// Declare a handler which writes a plain-text response
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map which holds the information that we want to send in the response.
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := serializer.ToJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w,
			"The server encountered a problem and could not process your request", http.StatusInternalServerError,
		)
	}
}
