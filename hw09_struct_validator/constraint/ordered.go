package constraint

import (
	"cmp"
	"errors"
	"reflect"
	"strconv"

	"github.com/kotoproger/go_hw/hw09structvalidator/core"
)

var (
	ErrInvalidValueLength = errors.New("invalid value length")
	ErrValueToLarge       = errors.New("value too large")
	ErrValueToSmall       = errors.New("value too small")
)

type orderedConstraint[T cmp.Ordered] struct {
	baseConstraint[T, T]
}

func NewMaxIntConstraint() core.ConstraintInterface {
	return newIntConstraint(
		func(source int64, target int64) bool {
			return source > target
		},
		"max",
		ErrValueToLarge,
		"max value %d, actual value %d: %w",
	)
}

func NewMinIntConstraint() core.ConstraintInterface {
	return newIntConstraint(
		func(source int64, target int64) bool {
			return source < target
		},
		"min",
		ErrValueToSmall,
		"min value %d, actual value %d: %w",
	)
}

func newIntConstraint(
	comp func(source int64, target int64) bool,
	name string,
	baseError error,
	message string,
) core.ConstraintInterface {
	return &orderedConstraint[int64]{
		baseConstraint[int64, int64]{
			kindTypes: []reflect.Kind{
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			},
			prepareValue: func(kind reflect.Value) int64 {
				return kind.Int()
			},
			prepareParam: func(params string) (value int64, err error) {
				return strconv.ParseInt(params, 10, 64)
			},
			compareValues:   comp,
			messageTemplate: message,
			name:            name,
			baseError:       baseError,
		},
	}
}

func NewMaxUintConstraint() core.ConstraintInterface {
	return newUintConstraint(
		func(source uint64, target uint64) bool {
			return source > target
		},
		"max",
		ErrValueToLarge,
		"max value %d, actual value %d: %w",
	)
}

func NewMinUintConstraint() core.ConstraintInterface {
	return newUintConstraint(
		func(source uint64, target uint64) bool {
			return source < target
		},
		"min",
		ErrValueToSmall,
		"min value %d, actual value %d: %w",
	)
}

func newUintConstraint(
	comp func(source uint64, target uint64) bool,
	name string,
	baseError error,
	message string,
) core.ConstraintInterface {
	return &orderedConstraint[uint64]{
		baseConstraint[uint64, uint64]{
			kindTypes: []reflect.Kind{
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			},
			prepareValue: func(kind reflect.Value) uint64 {
				return kind.Uint()
			},
			prepareParam: func(params string) (value uint64, err error) {
				return strconv.ParseUint(params, 10, 64)
			},
			compareValues:   comp,
			messageTemplate: message,
			name:            name,
			baseError:       baseError,
		},
	}
}

func newFloatConstraint(
	comp func(source float64, target float64) bool,
	name string,
	baseError error,
	message string,
) core.ConstraintInterface {
	return &orderedConstraint[float64]{
		baseConstraint[float64, float64]{
			kindTypes: []reflect.Kind{
				reflect.Float32, reflect.Float64,
			},
			prepareValue: func(kind reflect.Value) float64 {
				return kind.Float()
			},
			prepareParam: func(params string) (value float64, err error) {
				return strconv.ParseFloat(params, 64)
			},
			compareValues:   comp,
			messageTemplate: message,
			name:            name,
			baseError:       baseError,
		},
	}
}

func NewMaxFloatConstraint() core.ConstraintInterface {
	return newFloatConstraint(
		func(source float64, target float64) bool {
			return source > target
		},
		"max",
		ErrValueToLarge,
		"max value %d, actual value %d: %w",
	)
}

func NewMinFloatConstraint() core.ConstraintInterface {
	return newFloatConstraint(
		func(source float64, target float64) bool {
			return source < target
		},
		"min",
		ErrValueToSmall,
		"min value %d, actual value %d: %w",
	)
}

func NewLengthConstraint() core.ConstraintInterface {
	return &orderedConstraint[int]{
		baseConstraint[int, int]{
			kindTypes: []reflect.Kind{
				reflect.String, reflect.Slice, reflect.Chan, reflect.Array, reflect.Map,
			},
			prepareValue: func(value reflect.Value) int {
				return value.Len()
			},
			prepareParam: strconv.Atoi,
			compareValues: func(source int, target int) bool {
				return source != target
			},
			messageTemplate: "expected %d length but got %d: %w",
			name:            "len",
			baseError:       ErrInvalidValueLength,
		},
	}
}

func init() {
	core.RegisterConstraint(NewMaxIntConstraint())
	core.RegisterConstraint(NewMinIntConstraint())
	core.RegisterConstraint(NewMaxFloatConstraint())
	core.RegisterConstraint(NewMinFloatConstraint())
	core.RegisterConstraint(NewLengthConstraint())
	core.RegisterConstraint(NewMaxUintConstraint())
	core.RegisterConstraint(NewMinUintConstraint())
}
