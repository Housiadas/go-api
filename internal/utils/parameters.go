package utils

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func ReadIDParam(r *http.Request, paramName string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName(paramName), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New(fmt.Sprintf("invalid %v parameter", paramName))
	}
	return id, nil
}
