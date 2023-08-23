package main

import "log"

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

const version = "1.0.0"
