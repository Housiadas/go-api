package utils

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-api/internal/validator"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ReadIDParam(r *http.Request, paramName string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName(paramName), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New(fmt.Sprintf("invalid %v parameter", paramName))
	}
	return id, nil
}

// ReadString helper returns a string value from the query string, or the provided
// default value if no matching key could be found.
func ReadString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists this
	// will return the empty string "".
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	return s
}

// ReadCSV helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if csv == "" {
		return defaultValue
	}
	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

// ReadInt helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message in the provided Validator instance.
func ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Extract the value from the query string.
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	// Try to convert the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	// Otherwise, return the converted integer value.
	return i
}
