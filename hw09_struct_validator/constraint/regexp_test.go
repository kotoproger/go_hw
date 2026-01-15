package constraint

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/kotoproger/go_hw/hw09structvalidator/core"
	"github.com/stretchr/testify/assert"
)

func TestRegexpConstraintErrUnsupportedValueType(t *testing.T) {
	testCases := []struct {
		name       string
		value      interface{}
		validators []core.ConstraintInterface
	}{
		{
			name:  "bool",
			value: false,
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "complex64",
			value: complex64(123),
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "complex128",
			value: complex128(123),
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "array",
			value: [1]any{},
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "chan",
			value: make(chan any),
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "func",
			value: func() {},
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "map",
			value: map[any]any{},
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "slice",
			value: make([]any, 1),
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "struct",
			value: struct{}{},
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "int",
			value: 10,
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
			},
		},
		{
			name:  "float",
			value: 10.5,
			validators: []core.ConstraintInterface{
				NewRegexpConstraint(),
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

func TestRegexpConstraintCompileError(t *testing.T) {
	constraints, err := NewRegexpConstraint().Validate(reflect.ValueOf("qwe"), "asd(asd-]")

	var expectedErr core.ErrUnsupportedConstraintParams
	assert.Nil(t, constraints)
	assert.ErrorAs(t, err, &expectedErr)
}

func TestRegexpConstraintNotMatch(t *testing.T) {
	constraints, err := NewRegexpConstraint().Validate(reflect.ValueOf("qwe"), "`asd32`")

	assert.Nil(t, err)
	assert.NotNil(t, constraints)
	assert.ErrorIs(t, constraints[0], ErrDoesNotMatchPattern)
}

func TestRegexpConstraintMatch(t *testing.T) {
	constraints, err := NewRegexpConstraint().Validate(reflect.ValueOf("qwe"), "`.*`")

	assert.Nil(t, err)
	assert.NotNil(t, constraints)
	assert.ErrorIs(t, constraints[0], ErrDoesNotMatchPattern)
}
