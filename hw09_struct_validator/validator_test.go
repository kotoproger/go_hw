package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kotoproger/go_hw/hw09structvalidator/constraint"
	"github.com/kotoproger/go_hw/hw09structvalidator/core"
	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36|regexp:^[0-9]+$"`
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

	failure struct {
		String string `validate:"min:1|max:100"`
		Struct App    `validate:"len:5"`
	}
)

type expected struct {
	name string
	err  error
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr []expected
	}{
		{
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "1.0.0,3"},
			expectedErr: []expected{{name: "Version", err: constraint.ErrInvalidValueLength}},
		},
		{
			in:          App{Version: "1.0."},
			expectedErr: []expected{{name: "Version", err: constraint.ErrInvalidValueLength}},
		},
		{
			in:          Response{Code: 300, Body: "body"},
			expectedErr: []expected{{name: "Code", err: constraint.ErrValueIsNotAllowed}},
		},
		{
			in:          Response{Code: 200, Body: "body"},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456123456123456123456123456123456",
				Age:    22,
				Email:  "some@mail.ru",
				Role:   UserRole("admin"),
				Phones: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "1"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456123456123456123456123456123456",
				Age:    0,
				Email:  "some@mail.ru",
				Role:   UserRole("admin"),
				Phones: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "1"},
			},
			expectedErr: []expected{{name: "Age", err: constraint.ErrValueToSmall}},
		},
		{
			in: User{
				ID:     "123456123456123456123456123456123456",
				Age:    110,
				Email:  "some@mail.ru",
				Role:   UserRole("admin"),
				Phones: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "1"},
			},
			expectedErr: []expected{{name: "Age", err: constraint.ErrValueToLarge}},
		},
		{
			in: User{
				ID:     "123456123456123456123456123456123456",
				Age:    22,
				Email:  "some@mail.ru",
				Role:   UserRole("role"),
				Phones: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "1"},
			},
			expectedErr: []expected{{name: "Role", err: constraint.ErrValueIsNotAllowed}},
		},
		{
			in: User{
				ID:     "12345612345612345612345612345d12",
				Age:    100,
				Email:  "some!mail.ru",
				Role:   UserRole("role"),
				Phones: []string{"1", "2", "3", "4", "5", "6", "9", "0", "1"},
			},
			expectedErr: []expected{
				{name: "ID", err: constraint.ErrInvalidValueLength},
				{name: "ID", err: constraint.ErrDoesNotMatchPattern},
				{name: "Age", err: constraint.ErrValueToLarge},
				{name: "Email", err: constraint.ErrDoesNotMatchPattern},
				{name: "Role", err: constraint.ErrValueIsNotAllowed},
				{name: "Phones", err: constraint.ErrInvalidValueLength},
			},
		},
		{
			in: failure{},
			expectedErr: []expected{
				{name: "String", err: core.ErrUnsupportedValueType},
			},
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				errors, ok := err.(ValidationErrors)
				if ok {
					assert.Equal(t, 0, len(errors))
				} else {
					assert.Nil(t, err)
				}

			} else {
				errors, ok := err.(ValidationErrors)
				if !ok {
					assert.Equal(t, 1, len(tt.expectedErr), "expected more than 1 error")
					assert.ErrorIs(t, err, tt.expectedErr[0].err)
				} else {
					for index, e := range tt.expectedErr {
						assert.Equal(t, e.name, errors[index].Field)
						assert.ErrorIs(t, errors[index].Err, e.err)
					}
				}
			}
			_ = tt
		})
	}
}

func TestUnsupportedConstraintParams(t *testing.T) {
	testCases := []struct {
		name string
		in   any
	}{
		{
			name: "string length",
			in: struct {
				name string `validate:"len:adasd"`
			}{},
		},
		{
			name: "float length",
			in: struct {
				name string `validate:"len:5.45"`
			}{},
		},
		{
			name: "broken regexp",
			in: struct {
				name string `validate:"regexp:^[qwewqe("`
			}{},
		},
		{
			name: "params count",
			in: struct {
				name string `validate:"len:4:5"`
			}{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var expectedErr core.ErrUnsupportedConstraintParams
			err := Validate(tc.in)

			assert.ErrorAs(t, err, &expectedErr)
		})
	}
}

func TestUnknownConstratint(t *testing.T) {
	err := Validate(struct {
		unknownConstraint string `validate:"unknown:2121"`
	}{})
	var unknown core.ErrUnknownConstraint

	assert.ErrorAs(t, err, &unknown)
}
