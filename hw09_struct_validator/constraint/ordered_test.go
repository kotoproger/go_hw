package constraint

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kotoproger/go_hw/hw09structvalidator"
	"github.com/stretchr/testify/assert"
)

func TestMaxMinConstraintErrUnsupportedValueType(t *testing.T) { //nolint:funlen
	testCases := []struct {
		name       string
		value      interface{}
		validators []hw09structvalidator.ConstraintInterface
	}{
		{
			name:  "string",
			value: "string",
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "bool",
			value: false,
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewLengthConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "complex64",
			value: complex64(123),
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewLengthConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "complex128",
			value: complex128(123),
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewLengthConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "array",
			value: [1]any{},
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "chan",
			value: make(chan any),
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "func",
			value: func() {},
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewLengthConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "map",
			value: map[any]any{},
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "slice",
			value: make([]any, 1),
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "struct",
			value: struct{}{},
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewMinIntConstraint(),
				NewLengthConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "int",
			value: 10,
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
				NewLengthConstraint(),
			},
		},
		{
			name:  "float",
			value: 10.5,
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMinIntConstraint(),
				NewLengthConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, validator := range testCase.validators {
				t.Run(
					fmt.Sprintf("%s(%s)", validator.ConstraintName(), validator.ConstraintSubName()),
					func(t *testing.T) {
						constraints, err := validator.Validate(reflect.ValueOf(testCase.value), "6")

						assert.Nil(t, constraints)
						assert.Equal(t, hw09structvalidator.ErrUnsupportedValueType, err)
					},
				)
			}
		})
	}
}

func TestMaxMinConstraintErrUnsupportedParam(t *testing.T) {
	testCases := []struct {
		name       string
		value      interface{}
		param      string
		validators []hw09structvalidator.ConstraintInterface
	}{
		{
			name:  "int value - string parameter",
			value: 10,
			param: "int",
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "int value - float parameter",
			value: 10,
			param: "1.1",
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},

		{
			name:  "uint value - string parameter",
			value: uint(10),
			param: "int",
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},
		{
			name:  "uint value - float parameter",
			value: uint(10),
			param: "1.1",
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxIntConstraint(),
				NewMinIntConstraint(),
				NewMaxUintConstraint(),
				NewMinUintConstraint(),
			},
		},

		{
			name:  "float value - string parameter",
			value: 10.5,
			param: "int",
			validators: []hw09structvalidator.ConstraintInterface{
				NewMaxFloatConstraint(),
				NewMinFloatConstraint(),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, validator := range testCase.validators {
				t.Run(
					fmt.Sprintf("%s(%s)", validator.ConstraintName(), validator.ConstraintSubName()),
					func(t *testing.T) {
						constraints, err := validator.Validate(reflect.ValueOf(testCase.value), testCase.param)

						var expectedErr hw09structvalidator.ErrUnsupportedConstraintParams
						assert.Nil(t, constraints)
						assert.ErrorAs(t, err, &expectedErr)
					},
				)
			}
		})
	}
}

func TestMaxConstraintValid(t *testing.T) {
	testCases := []struct {
		name      string
		value     interface{}
		param     string
		validator hw09structvalidator.ConstraintInterface
	}{
		{
			name:      "int 5 = 5",
			value:     int(5),
			param:     "5",
			validator: NewMaxIntConstraint(),
		},
		{
			name:      "int 5 < 10",
			value:     int(5),
			param:     "10",
			validator: NewMaxIntConstraint(),
		},
		{
			name:      "int -5 < 10",
			value:     int(-5),
			param:     "10",
			validator: NewMaxIntConstraint(),
		},
		{
			name:      "int 0 < 10",
			value:     int(0),
			param:     "10",
			validator: NewMaxIntConstraint(),
		},
		{
			name:      "int -5 < -3",
			value:     int(-5),
			param:     "-3",
			validator: NewMaxIntConstraint(),
		},

		{
			name:      "float 5.5 = 5.5",
			value:     float32(5.5),
			param:     "5.5",
			validator: NewMaxFloatConstraint(),
		},
		{
			name:      "float 5.5 < 5.6",
			value:     float32(5.5),
			param:     "5.6",
			validator: NewMaxFloatConstraint(),
		},
		{
			name:      "float -5.5 < 12",
			value:     float32(-5.5),
			param:     "12",
			validator: NewMaxFloatConstraint(),
		},
		{
			name:      "float 0 < 10",
			value:     float32(0),
			param:     "10",
			validator: NewMaxFloatConstraint(),
		},
		{
			name:      "float -5.5 < -5.4",
			value:     float32(-5.5),
			param:     "-5.4",
			validator: NewMaxFloatConstraint(),
		},

		{
			name:      "uint 5 = 5",
			value:     uint(5),
			param:     "5",
			validator: NewMaxUintConstraint(),
		},
		{
			name:      "uint 5 < 10",
			value:     uint(5),
			param:     "10",
			validator: NewMaxUintConstraint(),
		},
		{
			name:      "uint 0 < 10",
			value:     uint(0),
			param:     "10",
			validator: NewMaxUintConstraint(),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			constraints, err := testCase.validator.Validate(reflect.ValueOf(testCase.value), testCase.param)

			assert.Nil(t, err)
			assert.Nil(t, constraints)
		})
	}
}

func TestMaxConstraintInvalid(t *testing.T) {
	testCases := []struct {
		name      string
		value     interface{}
		param     string
		validator hw09structvalidator.ConstraintInterface
	}{
		{
			name:      "int 5 !< 4",
			value:     int(5),
			param:     "4",
			validator: NewMaxIntConstraint(),
		},
		{
			name:      "int -5 !< -10",
			value:     int(-5),
			param:     "-10",
			validator: NewMaxIntConstraint(),
		},
		{
			name:      "int 0 !< -1",
			value:     int(0),
			param:     "-1",
			validator: NewMaxIntConstraint(),
		},

		{
			name:      "float 5.5 !< 5.4",
			value:     float32(5.5),
			param:     "5.4",
			validator: NewMaxFloatConstraint(),
		},
		{
			name:      "float -5.5 !< -12",
			value:     float32(-5.5),
			param:     "-12",
			validator: NewMaxFloatConstraint(),
		},

		{
			name:      "float -5.5 < -5.6",
			value:     float32(-5.5),
			param:     "-5.6",
			validator: NewMaxFloatConstraint(),
		},

		{
			name:      "uint 5 !< 4",
			value:     uint(5),
			param:     "4",
			validator: NewMaxUintConstraint(),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			constraints, err := testCase.validator.Validate(reflect.ValueOf(testCase.value), testCase.param)

			assert.Nil(t, err)
			assert.NotNil(t, constraints)
			assert.ErrorIs(t, constraints[0], ErrValueToLarge)
		})
	}
}

func TestMinConstraintValid(t *testing.T) {
	testCases := []struct {
		name      string
		value     interface{}
		param     string
		validator hw09structvalidator.ConstraintInterface
	}{
		{
			name:      "int 5 = 5",
			value:     int(5),
			param:     "5",
			validator: NewMinIntConstraint(),
		},
		{
			name:      "int 5 > 4",
			value:     int(5),
			param:     "4",
			validator: NewMinIntConstraint(),
		},
		{
			name:      "int -5 > -10",
			value:     int(-5),
			param:     "-10",
			validator: NewMinIntConstraint(),
		},
		{
			name:      "int 0 > -10",
			value:     int(0),
			param:     "-10",
			validator: NewMinIntConstraint(),
		},
		{
			name:      "int -5 > -6",
			value:     int(-5),
			param:     "-6",
			validator: NewMinIntConstraint(),
		},

		{
			name:      "float 5.5 = 5.5",
			value:     float32(5.5),
			param:     "5.5",
			validator: NewMinFloatConstraint(),
		},
		{
			name:      "float 5.5 > 5.4",
			value:     float32(5.5),
			param:     "5.4",
			validator: NewMinFloatConstraint(),
		},
		{
			name:      "float -5.5 > -12",
			value:     float32(-5.5),
			param:     "-12",
			validator: NewMinFloatConstraint(),
		},
		{
			name:      "float 0 > -10",
			value:     float32(0),
			param:     "-10",
			validator: NewMinFloatConstraint(),
		},
		{
			name:      "float -5.5 > -5.6",
			value:     float32(-5.5),
			param:     "-5.6",
			validator: NewMinFloatConstraint(),
		},

		{
			name:      "uint 5 = 5",
			value:     uint(5),
			param:     "5",
			validator: NewMinUintConstraint(),
		},
		{
			name:      "uint 5 > 4",
			value:     uint(5),
			param:     "4",
			validator: NewMinUintConstraint(),
		},
		{
			name:      "uint 1 > 0",
			value:     uint(1),
			param:     "0",
			validator: NewMinUintConstraint(),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			constraints, err := testCase.validator.Validate(reflect.ValueOf(testCase.value), testCase.param)

			assert.Nil(t, err)
			assert.Nil(t, constraints)
		})
	}
}

func TestMinConstraintInvalid(t *testing.T) {
	testCases := []struct {
		name      string
		value     interface{}
		param     string
		validator hw09structvalidator.ConstraintInterface
	}{
		{
			name:      "int 5 !> 6",
			value:     int(5),
			param:     "6",
			validator: NewMinIntConstraint(),
		},
		{
			name:      "int -5 !> -4",
			value:     int(-5),
			param:     "-4",
			validator: NewMinIntConstraint(),
		},
		{
			name:      "int 0 !> 1",
			value:     int(0),
			param:     "1",
			validator: NewMinIntConstraint(),
		},

		{
			name:      "float 5.5 !> 5.6",
			value:     float32(5.5),
			param:     "5.6",
			validator: NewMinFloatConstraint(),
		},
		{
			name:      "float -5.5 !> -1",
			value:     float32(-5.5),
			param:     "-1",
			validator: NewMinFloatConstraint(),
		},

		{
			name:      "float -5.5 > -5.4",
			value:     float32(-5.5),
			param:     "-5.4",
			validator: NewMinFloatConstraint(),
		},

		{
			name:      "uint 5 !> 6",
			value:     uint(5),
			param:     "6",
			validator: NewMinUintConstraint(),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			constraints, err := testCase.validator.Validate(reflect.ValueOf(testCase.value), testCase.param)

			assert.Nil(t, err)
			assert.NotNil(t, constraints)
			assert.ErrorIs(t, constraints[0], ErrValueToSmall)
		})
	}
}

func TestLengthConstraintUnsupportedParam(t *testing.T) {
	testCases := []struct {
		name   string
		params string
	}{
		{
			name:   "some text as params",
			params: "some text",
		},
		{
			name:   "empty string",
			params: "",
		},
		{
			name:   "float falue",
			params: "23.3231",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			constraints, err := NewLengthConstraint().Validate(reflect.ValueOf("asdads"), tc.params)

			var expectedErr hw09structvalidator.ErrUnsupportedConstraintParams
			assert.ErrorAs(t, err, &expectedErr)
			assert.Nil(t, constraints)
		})
	}
}

func TestLengthConstraintValid(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "string",
			value: "string",
		},
		{
			name:  "slice",
			value: make([]any, 6),
		},
		{
			name:  "chan",
			value: make([]chan int, 6),
		},
		{
			name:  "array",
			value: [6]any{},
		},
		{
			name: "map",
			value: map[any]any{
				1: 2,
				3: 4,
				5: 6,
				7: 8,
				9: 0,
				2: 1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			constraints, err := NewLengthConstraint().Validate(reflect.ValueOf(tc.value), "6")

			assert.Nil(t, err)
			assert.Nil(t, constraints)
		})
	}
}

func TestLengthConstraintInvalid(t *testing.T) {
	testCases := []struct {
		name  string
		value any
	}{
		{
			name:  "short string",
			value: "strin",
		},
		{
			name:  "long string",
			value: "strinigs",
		},
		{
			name:  "short slice",
			value: make([]any, 5),
		},
		{
			name:  "long slice",
			value: make([]any, 7),
		},
		{
			name:  "short chan",
			value: make([]chan int, 5),
		},
		{
			name:  "long chan",
			value: make([]chan int, 7),
		},
		{
			name:  "short array",
			value: [5]any{},
		},
		{
			name:  "long array",
			value: [7]any{},
		},
		{
			name: "sort map",
			value: map[any]any{
				1: 2,
				3: 4,
				5: 6,
				7: 8,
				9: 0,
			},
		},
		{
			name: "long map",
			value: map[any]any{
				1: 2,
				3: 4,
				5: 6,
				7: 8,
				9: 0,
				2: 1,
				4: 1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			constraints, err := NewLengthConstraint().Validate(reflect.ValueOf(tc.value), "6")

			assert.Nil(t, err)
			assert.NotNil(t, constraints)
			assert.ErrorIs(t, constraints[0], ErrInvalidValueLength)
		})
	}
}
