package validator

import (
	"go-api/internal/models"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	type args struct {
		v     *Validator
		email string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Empty email input",
			args: args{
				v:     New(),
				email: "",
			},
			expected: false,
		},
		{
			name: "Invalid email input",
			args: args{
				v:     New(),
				email: "test",
			},
			expected: false,
		},
		{
			name: "Valid email input",
			args: args{
				v:     New(),
				email: "test@test.com",
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateEmail(tt.args.v, tt.args.email)
			if tt.args.v.Valid() != tt.expected {
				t.Errorf("got %v expected %v", tt.args.v.Valid(), tt.expected)
			}
		})
	}
}

func TestValidatePasswordPlaintext(t *testing.T) {
	type args struct {
		v        *Validator
		password string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Empty password",
			args: args{
				v:        New(),
				password: "",
			},
			expected: false,
		},
		{
			name: "Password under 8 bytes",
			args: args{
				v:        New(),
				password: "123",
			},
			expected: false,
		},
		{
			name: "Valid Password",
			args: args{
				v:        New(),
				password: "12345678",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidatePasswordPlaintext(tt.args.v, tt.args.password)
			if tt.args.v.Valid() != tt.expected {
				t.Errorf("got %v expected %v", tt.args.v.Valid(), tt.expected)
			}
		})
	}
}

func TestValidateUserEmptyName(t *testing.T) {
	v := New()
	user := &models.User{
		Name:  "",
		Email: "test@test.com",
	}
	user.Password.Set("123456789")
	ValidateUser(v, user)
	if v.Valid() != false {
		t.Errorf("got %v expected %v", v.Valid(), false)
	}
}

func TestValidateUserInvalidPassword(t *testing.T) {
	v := New()
	user := &models.User{
		Name:  "chris housi",
		Email: "test@test.com",
	}
	user.Password.Set("123")
	ValidateUser(v, user)
	if v.Valid() != false {
		t.Errorf("got %v expected %v", v.Valid(), false)
	}
}

func TestValidateUser(t *testing.T) {
	v := New()
	user := &models.User{
		Name:  "top gun",
		Email: "test@test.com",
	}
	user.Password.Set("123456789")
	ValidateUser(v, user)
	if v.Valid() != true {
		t.Errorf("got %v expected %v", v.Valid(), true)
	}
}
