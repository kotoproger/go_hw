package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kotoproger/go_hw/hw09structvalidator/constraint"
	"github.com/stretchr/testify/assert"
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
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "1.0.0,3"},
			expectedErr: constraint.ErrInvalidValueLength,
		},
		{
			in:          App{Version: "1.0."},
			expectedErr: constraint.ErrInvalidValueLength,
		},
		{
			in:          Response{Code: 300, Body: "body"},
			expectedErr: constraint.ErrValueIsNotAllowed,
		},
		{
			in:          Response{Code: 200, Body: "body"},
			expectedErr: nil,
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				assert.Equal(t, 0, len(err))
			} else {
				assert.ErrorIs(t, err[0].Err, tt.expectedErr)
			}
			_ = tt
		})
	}
}
