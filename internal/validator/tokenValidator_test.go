package validator

import "testing"

func TestValidateTokenPlaintext(t *testing.T) {
	type args struct {
		v              *Validator
		tokenPlaintext string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Empty token",
			args: args{
				v:              New(),
				tokenPlaintext: "",
			},
			expected: false,
		},
		{
			name: "Token less than 26 bytes",
			args: args{
				v:              New(),
				tokenPlaintext: "OO2VLM3C2L6",
			},
			expected: false,
		},
		{
			name: "Valid token",
			args: args{
				v:              New(),
				tokenPlaintext: "OO2VLM3C2L6UO75DXMKNEWJYFA",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateTokenPlaintext(tt.args.v, tt.args.tokenPlaintext)
			if tt.args.v.Valid() != tt.expected {
				t.Errorf("got %v expected %v", tt.args.v.Valid(), tt.expected)
			}
		})
	}
}
