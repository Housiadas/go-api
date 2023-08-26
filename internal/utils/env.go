package utils

import (
	"os"
	"strconv"
)

func GetenvInt(key string) (int, error) {
	s := os.Getenv(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetenvBool(key string) (bool, error) {
	s := os.Getenv(key)
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
