package constraint

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/kotoproger/go_hw/hw09structvalidator/core"
)

type baseConstraint[T any, S any] struct {
	kindTypes       []reflect.Kind
	prepareValue    func(value reflect.Value) T
	prepareParam    func(params string) (S, error)
	compareValues   func(source T, target S) bool
	messageTemplate string
	name            string
	baseError       error
}

func (g *baseConstraint[T, S]) ConstraintName() string {
	return g.name
}

func (g *baseConstraint[T, S]) KindTypes() []reflect.Kind {
	return g.kindTypes
}

func (g *baseConstraint[T, S]) Validate(value reflect.Value, params string) (constraints []error, err error) {
	if !slices.Contains(g.kindTypes, value.Type().Kind()) {
		err = core.ErrUnsupportedValueType
		return
	}

	preparedValue := g.prepareValue(value)
	preparedExpects, prepareErr := g.prepareParam(params)
	if prepareErr != nil {
		err = core.ErrUnsupportedConstraintParams(prepareErr)

		return
	}

	if g.compareValues(preparedValue, preparedExpects) {
		constraints = append(constraints, fmt.Errorf(g.messageTemplate, preparedValue, preparedExpects, g.baseError))
	}

	return
}
