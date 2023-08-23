package main

import (
	"encoding/json"
	"net/http"
)

// Declare a handler which writes a plain-text response
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Create a map which holds the information that we want to send in the response.
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}

	// Append a newline to the JSON. This is just a small nicety to make it easier to view in terminal applications.
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		return
	}
}
