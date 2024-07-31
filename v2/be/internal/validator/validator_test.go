package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMinString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		length  int
		field   string
		message string
		wantErr bool
	}{
		{
			name:    "String length equal to minimum",
			s:       "test",
			length:  4,
			field:   "testField",
			message: "String too short",
			wantErr: false,
		},
		{
			name:    "String length greater than minimum",
			s:       "testing",
			length:  5,
			field:   "testField",
			message: "String too short",
			wantErr: false,
		},
		{
			name:    "Empty string",
			s:       "",
			length:  1,
			field:   "testField",
			message: "String too short",
			wantErr: true,
		},
		{
			name:    "Unicode string", // modify based on preference
			s:       "こんにちは",
			length:  4,
			field:   "testField",
			message: "String too short",
			wantErr: false,
		},
		{
			name:    "String with spaces",
			s:       "  a  ",
			length:  4,
			field:   "testField",
			message: "String too short",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.MinString(tt.s, tt.length, tt.field, tt.message)

			if tt.wantErr {
				require.False(t, v.Valid())
				require.NotEmpty(t, v.Errors())
				require.Equal(t, tt.message, v.Errors()[tt.field])
			} else {
				require.True(t, v.Valid())
				require.Empty(t, v.Errors())
			}
		})
	}
}

func TestRequiredString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		field   string
		message string
		wantErr bool
	}{
		{
			name:    "Non-empty string",
			s:       "hello",
			field:   "testField",
			message: "Field is required",
			wantErr: false,
		},
		{
			name:    "Empty string",
			s:       "",
			field:   "testField",
			message: "Field is required",
			wantErr: true,
		},
		{
			name:    "String with only spaces",
			s:       "   ",
			field:   "testField",
			message: "Field is required",
			wantErr: true,
		},
		{
			name:    "String with special characters",
			s:       "!@#$%^&*()",
			field:   "testField",
			message: "Field is required",
			wantErr: false,
		},
		{
			name:    "Unicode string",
			s:       "こんにちは",
			field:   "testField",
			message: "Field is required",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.RequiredString(tt.s, tt.field, tt.message)

			if tt.wantErr {
				require.False(t, v.Valid())
				require.NotEmpty(t, v.Errors())
				require.Equal(t, tt.message, v.Errors()[tt.field])
			} else {
				require.True(t, v.Valid())
				require.Empty(t, v.Errors())
			}
		})
	}
}

func TestAddError(t *testing.T) {
	tests := []struct {
		name           string
		existingErrors map[string]string
		field          string
		message        string
		expectedErrors map[string]string
	}{
		{
			name:           "Add error to empty validator",
			existingErrors: map[string]string{},
			field:          "username",
			message:        "Username is required",
			expectedErrors: map[string]string{"username": "Username is required"},
		},
		{
			name:           "Add error to validator with existing errors",
			existingErrors: map[string]string{"email": "Invalid email format"},
			field:          "password",
			message:        "Password is too short",
			expectedErrors: map[string]string{
				"email":    "Invalid email format",
				"password": "Password is too short",
			},
		},
		{
			name:           "Add error for field that already has an error",
			existingErrors: map[string]string{"username": "Username is required"},
			field:          "username",
			message:        "Username must be unique",
			expectedErrors: map[string]string{"username": "Username is required"},
		},
		{
			name:           "Add error with empty field name",
			existingErrors: map[string]string{},
			field:          "",
			message:        "Generic error",
			expectedErrors: map[string]string{"": "Generic error"},
		},
		{
			name:           "Add error with empty message",
			existingErrors: map[string]string{},
			field:          "age",
			message:        "",
			expectedErrors: map[string]string{"age": ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{errs: tt.existingErrors}
			v.AddError(tt.field, tt.message)

			require.Equal(t, tt.expectedErrors, v.errs)
			require.Len(t, v.errs, len(tt.expectedErrors))

			for field, expectedMsg := range tt.expectedErrors {
				actualMsg, exists := v.errs[field]
				require.True(t, exists)
				require.Equal(t, expectedMsg, actualMsg)
			}
		})
	}
}
