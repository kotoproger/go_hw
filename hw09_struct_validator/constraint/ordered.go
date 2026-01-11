package constraint

import (
	"cmp"
	"errors"
	"reflect"
	"strconv"

	"github.com/kotoproger/go_hw/hw09structvalidator"
)

var (
	ErrInvalidValueLength = errors.New("invalid value length")
	ErrValueToLarge       = errors.New("value too large")
	ErrValueToSmall       = errors.New("value too small")
)

type orderedConstraint[T cmp.Ordered] struct {
	baseConstraint[T, T]
}

func NewMaxIntConstraint() hw09structvalidator.ConstraintInterface {
	return newIntConstraint(
		func(source int64, target int64) bool {
			return source > target
		},
		"max",
		"int",
		ErrValueToLarge,
		"max value %d, actual value %d: %w",
	)
}

func NewMinIntConstraint() hw09structvalidator.ConstraintInterface {
	return newIntConstraint(
		func(source int64, target int64) bool {
			return source < target
		},
		"min",
		"int",
		ErrValueToSmall,
		"min value %d, actual value %d: %w",
	)
}

func newIntConstraint(
	comp func(source int64, target int64) bool,
	name string,
	subname string,
	baseError error,
	message string,
) hw09structvalidator.ConstraintInterface {
	return orderedConstraint[int64]{
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
			subname:         subname,
			baseError:       baseError,
		},
	}
}

func NewMaxUintConstraint() hw09structvalidator.ConstraintInterface {
	return newUintConstraint(
		func(source uint64, target uint64) bool {
			return source > target
		},
		"max",
		"uint",
		ErrValueToLarge,
		"max value %d, actual value %d: %w",
	)
}

func NewMinUintConstraint() hw09structvalidator.ConstraintInterface {
	return newUintConstraint(
		func(source uint64, target uint64) bool {
			return source < target
		},
		"min",
		"uint",
		ErrValueToSmall,
		"min value %d, actual value %d: %w",
	)
}

func newUintConstraint(
	comp func(source uint64, target uint64) bool,
	name string,
	subname string,
	baseError error,
	message string,
) hw09structvalidator.ConstraintInterface {
	return orderedConstraint[uint64]{
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
			subname:         subname,
			baseError:       baseError,
		},
	}
}

func newFloatConstraint(
	comp func(source float64, target float64) bool,
	name string,
	subname string,
	baseError error,
	message string,
) hw09structvalidator.ConstraintInterface {
	return orderedConstraint[float64]{
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
			subname:         subname,
			baseError:       baseError,
		},
	}
}

func NewMaxFloatConstraint() hw09structvalidator.ConstraintInterface {
	return newFloatConstraint(
		func(source float64, target float64) bool {
			return source > target
		},
		"max",
		"float",
		ErrValueToLarge,
		"max value %d, actual value %d: %w",
	)
}

func NewMinFloatConstraint() hw09structvalidator.ConstraintInterface {
	return newFloatConstraint(
		func(source float64, target float64) bool {
			return source < target
		},
		"min",
		"float",
		ErrValueToSmall,
		"min value %d, actual value %d: %w",
	)
}

func NewLengthConstraint() hw09structvalidator.ConstraintInterface {
	return orderedConstraint[int]{
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
			subname:         "int",
			baseError:       ErrInvalidValueLength,
		},
	}
}

func init() {
	hw09structvalidator.RegisterConstraint(NewMaxIntConstraint())
	hw09structvalidator.RegisterConstraint(NewMinIntConstraint())
	hw09structvalidator.RegisterConstraint(NewMaxFloatConstraint())
	hw09structvalidator.RegisterConstraint(NewMinFloatConstraint())
	hw09structvalidator.RegisterConstraint(NewLengthConstraint())
	hw09structvalidator.RegisterConstraint(NewMaxUintConstraint())
	hw09structvalidator.RegisterConstraint(NewMinUintConstraint())
}
