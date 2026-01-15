package core

import (
	"errors"
	"fmt"
	"reflect"
)

var Constraints = make(map[string]map[reflect.Kind]ConstraintInterface)

type (
	ErrUnsupportedConstraintParams error
	ErrUnknownConstraint           error
	ErrConstraint                  error
)

var ErrUnsupportedValueType = errors.New("unsupported value type")

// ErrUnsupportedConstraintParams = errors.New("unsupported constraint params").

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err.Error())
}

type ConstraintInterface interface {
	Validate(value reflect.Value, params string) ([]error, error)
	ConstraintName() string
	KindTypes() []reflect.Kind
}

func RegisterConstraint(constraint ConstraintInterface) {
	_, ok := Constraints[constraint.ConstraintName()]
	if !ok {
		Constraints[constraint.ConstraintName()] = make(map[reflect.Kind]ConstraintInterface)
	}
	for _, kindType := range constraint.KindTypes() {
		_, ok = Constraints[constraint.ConstraintName()][kindType]
		if ok {
			panic("duplicate constraint name: " + constraint.ConstraintName() + " fore type " + kindType.String())
		}
		Constraints[constraint.ConstraintName()][kindType] = constraint
	}
}
