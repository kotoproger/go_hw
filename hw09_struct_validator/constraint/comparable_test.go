package constraint

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kotoproger/go_hw/hw09structvalidator/core"
	"github.com/stretchr/testify/assert"
)

func TestComparableUnsupportedValueType(t *testing.T) {
	testCases := []struct {
		name       string
		value      interface{}
		validators []core.ConstraintInterface
	}{
		{
			name:  "bool",
			value: false,
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "complex64",
			value: complex64(123),
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "complex128",
			value: complex128(123),
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "array",
			value: [1]any{},
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "chan",
			value: make(chan any),
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "func",
			value: func() {},
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "map",
			value: map[any]any{},
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "slice",
			value: make([]any, 1),
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "struct",
			value: struct{}{},
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "string",
			value: "asdsd",
			validators: []core.ConstraintInterface{
				NewInFloatConstraint(),
				NewInIntConstraint(),
			},
		},
		{
			name:  "int",
			value: 10,
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInFloatConstraint(),
			},
		},
		{
			name:  "float",
			value: 10.5,
			validators: []core.ConstraintInterface{
				NewInStringConstraint(),
				NewInIntConstraint(),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, validator := range testCase.validators {
				t.Run(
					fmt.Sprintf("%s(%s)", validator.ConstraintName(), validator.KindTypes()),
					func(t *testing.T) {
						constraints, err := validator.Validate(reflect.ValueOf(testCase.value), "6")

						assert.Nil(t, constraints)
						assert.Equal(t, core.ErrUnsupportedValueType, err)
					},
				)
			}
		})
	}
}

func TestComparableConstraintErrUnsupportedParam(t *testing.T) {
	testCases := []struct {
		name       string
		value      interface{}
		validators []core.ConstraintInterface
		params     string
	}{
		{
			name:  "int(string)",
			value: 10,
			validators: []core.ConstraintInterface{
				NewInIntConstraint(),
			},
			params: `6,sdasd`,
		},
		{
			name:  "int(float)",
			value: 10,
			validators: []core.ConstraintInterface{
				NewInIntConstraint(),
			},
			params: `6,6.5`,
		},
		{
			name:  "float(string)",
			value: 10.5,
			validators: []core.ConstraintInterface{
				NewInFloatConstraint(),
			},
			params: `6,6.5,asdas`,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, validator := range testCase.validators {
				t.Run(
					fmt.Sprintf("%s(%s)", validator.ConstraintName(), validator.KindTypes()),
					func(t *testing.T) {
						constraints, err := validator.Validate(reflect.ValueOf(testCase.value), testCase.params)

						var expectedErr core.ErrUnsupportedConstraintParams
						assert.Nil(t, constraints)
						assert.ErrorAs(t, err, &expectedErr)
					},
				)
			}
		})
	}
}

func TestComparableConstraintValid(t *testing.T) {
	testCases := []struct {
		name      string
		value     interface{}
		validator core.ConstraintInterface
		params    string
	}{
		{
			name:      "int alone value",
			value:     10,
			validator: NewInIntConstraint(),
			params:    `10`,
		},
		{
			name:      "int first value",
			value:     10,
			validator: NewInIntConstraint(),
			params:    `10,123`,
		},
		{
			name:      "int last value",
			value:     10,
			validator: NewInIntConstraint(),
			params:    `123,23,10`,
		},
		{
			name:      "int zero value",
			value:     0,
			validator: NewInIntConstraint(),
			params:    `123,0,10`,
		},
		{
			name:      "int negative value",
			value:     -30,
			validator: NewInIntConstraint(),
			params:    `123,0,10,-30`,
		},

		{
			name:      "float",
			value:     10.65,
			validator: NewInFloatConstraint(),
			params:    `10.65`,
		},
		{
			name:      "float as int",
			value:     10.0,
			validator: NewInFloatConstraint(),
			params:    `10,123`,
		},
		{
			name:      "more than 1 zero",
			value:     10.0005,
			validator: NewInFloatConstraint(),
			params:    `123,23,10.0005`,
		},

		{
			name:      "string",
			value:     "asdsd",
			validator: NewInStringConstraint(),
			params:    `asdsd`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			constraints, err := testCase.validator.Validate(reflect.ValueOf(testCase.value), testCase.params)

			assert.Nil(t, err)
			assert.Nil(t, constraints)
		})
	}
}

func TestComparableConstraintInvalid(t *testing.T) {
	testCases := []struct {
		name      string
		value     interface{}
		validator core.ConstraintInterface
		params    string
	}{
		{
			name:      "int alone value",
			value:     10,
			validator: NewInIntConstraint(),
			params:    `11`,
		},
		{
			name:      "int first value",
			value:     10,
			validator: NewInIntConstraint(),
			params:    `22,123`,
		},
		{
			name:      "int last value",
			value:     10,
			validator: NewInIntConstraint(),
			params:    `123,23`,
		},
		{
			name:      "int zero value",
			value:     0,
			validator: NewInIntConstraint(),
			params:    `123,-10,10`,
		},

		{
			name:      "float",
			value:     10.65,
			validator: NewInFloatConstraint(),
			params:    `10.6`,
		},
		{
			name:      "float as int",
			value:     10.0,
			validator: NewInFloatConstraint(),
			params:    `10.123`,
		},
		{
			name:      "more than 1 zero",
			value:     10.0005,
			validator: NewInFloatConstraint(),
			params:    `123,23,10`,
		},

		{
			name:      "string",
			value:     "asdsd",
			validator: NewInStringConstraint(),
			params:    `asdsdqwe`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			constraints, err := testCase.validator.Validate(reflect.ValueOf(testCase.value), testCase.params)

			assert.Nil(t, err)
			assert.NotNil(t, constraints)
			assert.ErrorIs(t, constraints[0], ErrValueIsNotAllowed)
		})
	}
}
