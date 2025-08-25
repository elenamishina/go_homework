package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789112345678901234567890123456",
				Name:   "Elena",
				Age:    18,
				Email:  "elena@gmail.com",
				Role:   "stuff",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "some response",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "12345678911234567890",
				Name:   "Elena",
				Age:    17,
				Email:  "elenagmail.com",
				Role:   "stuf",
				Phones: []string{"1234567890"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrValidateLength},
				{Field: "Age", Err: ErrValidateMin},
				{Field: "Email", Err: ErrValidateRegexp},
				{Field: "Phones", Err: ErrValidateLength},
				{Field: "Role", Err: ErrValidateIn},
			},
		},
		{
			in: User{
				ID:     "123456789112345678901234567890123456",
				Name:   "Elena",
				Age:    65,
				Email:  "elena@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: ErrValidateMax},
			},
		},
		{
			in: App{
				Version: "v1.0.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrValidateLength},
			},
		},
		{
			in: Response{
				Code: 0,
				Body: "some response",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrValidateIn},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)

			checkValidationError(t, err, tt.expectedErr)
		})
	}
}

func checkValidationError(t *testing.T, err error, expectedErr error) {
	t.Helper()

	if expectedErr == nil {
		if err != nil {
			t.Errorf("Validate() error - %v, expected - no error", err)
		}
		return
	}
	if err == nil {
		t.Fatalf("Validate() - no error, expected - %v", expectedErr)
	}

	var expectValidErr ValidationErrors
	if errors.As(expectedErr, &expectValidErr) {
		var validateErr ValidationErrors
		if !errors.As(err, &validateErr) {
			t.Fatalf("Expected ValidationErrors type, got %v", err)
		}
		if len(validateErr) != len(expectValidErr) {
			t.Fatalf("Validate() num error - %d, expected - %d", len(validateErr), len(expectValidErr))
		}
		for _, wantErr := range expectValidErr {
			found := false
			for _, gotErr := range validateErr {
				if gotErr.Field == wantErr.Field && errors.Is(gotErr.Err, wantErr.Err) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Validate() error - not found, expected - %v", wantErr)
			}
		}
		return
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Validate() error - %v, expected - %v", err, expectedErr)
	}
}
