package utils

import "testing"

func TestGetenvBoolMissingKey(t *testing.T) {
	_, err := GetenvBool("dummy")
	if err == nil {
		t.Errorf("GetenvBool() expected error")
	}
}

func TestGetenvBool(t *testing.T) {
	t.Setenv("TEST_BOOL", "TRUE")
	got, err := GetenvBool("TEST_BOOL")
	if err != nil {
		t.Errorf("GetenvBool() error missing key")
		return
	}
	if got != true {
		t.Errorf("GetenvBool() got = %v, expected %v", got, true)
	}
}

func TestGetenvIntMissingKey(t *testing.T) {
	_, err := GetenvInt("dummy")
	if err == nil {
		t.Errorf("GetenvInt() expected error")
	}
}

func TestGetenvInt(t *testing.T) {
	t.Setenv("TEST_INT", "100")
	got, err := GetenvInt("TEST_INT")
	if err != nil {
		t.Errorf("GetenvInt() error missing key")
		return
	}
	if got != 100 {
		t.Errorf("GetenvInt() got = %v, expected %v", got, 100)
	}
}
