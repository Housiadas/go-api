package main

import (
	"fmt"
	"net/http"
)

// important to realize that our middleware will only recover panics that happen in
// the same goroutine that executed the recoverPanic() middleware.
// If, for example, you have a handler which spins up another goroutine (e.g. to do some
// background processing), then any panics that happen in the background goroutine will not
// be recovered — not by the recoverPanic() middleware… and not by the panic recovery
// built into http.Server . These panics will cause your application to exit and bring down the server.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not.
			err := recover()
			if err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type interface{}, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
