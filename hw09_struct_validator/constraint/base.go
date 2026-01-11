package constraint

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/kotoproger/go_hw/hw09structvalidator"
)

type baseConstraint[T any, S any] struct {
	kindTypes       []reflect.Kind
	prepareValue    func(value reflect.Value) T
	prepareParam    func(params string) (S, error)
	compareValues   func(source T, target S) bool
	messageTemplate string
	name            string
	subname         string
	baseError       error
}

func (g baseConstraint[T, S]) ConstraintName() string {
	return g.name
}

func (g baseConstraint[T, S]) ConstraintSubName() string {
	return g.subname
}

func (g baseConstraint[T, S]) Validate(value reflect.Value, params string) (constraints []error, err error) {
	if !slices.Contains(g.kindTypes, value.Type().Kind()) {
		err = hw09structvalidator.ErrUnsupportedValueType
		return
	}

	preparedValue := g.prepareValue(value)
	preparedExpects, prepareErr := g.prepareParam(params)
	if prepareErr != nil {
		err = hw09structvalidator.ErrUnsupportedConstraintParams(prepareErr)

		return
	}

	if g.compareValues(preparedValue, preparedExpects) {
		constraints = append(constraints, fmt.Errorf(g.messageTemplate, preparedValue, preparedExpects, g.baseError))
	}

	return
}
