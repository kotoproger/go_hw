package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

var constraints = make(map[string]map[string]ConstraintInterface)

type ErrUnsupportedConstraintParams error

var ErrUnsupportedValueType = errors.New("unsupported value type")

// ErrUnsupportedConstraintParams = errors.New("unsupported constraint params").

type ValidationError struct {
	Field string
	Err   error
}

type ConstraintInterface interface {
	Validate(value reflect.Value, params string) ([]error, error)
	ConstraintName() string
	ConstraintSubName() string
}

func RegisterConstraint(constraint ConstraintInterface) {
	_, ok := constraints[constraint.ConstraintName()]
	if !ok {
		constraints[constraint.ConstraintName()] = make(map[string]ConstraintInterface)
	}
	_, ok = constraints[constraint.ConstraintName()][constraint.ConstraintSubName()]
	if ok {
		panic(fmt.Sprintf("constraint %s(%s) already exists", constraint.ConstraintName(), constraint.ConstraintSubName()))
	}
	constraints[constraint.ConstraintName()][constraint.ConstraintSubName()] = constraint
}
