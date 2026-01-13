package constraint

import (
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/kotoproger/go_hw/hw09structvalidator/core"
)

var ErrValueIsNotAllowed core.ErrConstraint = errors.New("value is not allowed")

type comparableConstraint[T comparable] struct {
	baseConstraint[T, []T]
}

func NewInStringConstraint() core.ConstraintInterface {
	return &comparableConstraint[string]{
		baseConstraint[string, []string]{
			kindTypes: []reflect.Kind{reflect.String},
			prepareValue: func(value reflect.Value) string {
				return value.String()
			},
			prepareParam: func(params string) ([]string, error) {
				return strings.Split(params, ","), nil
			},
			compareValues: func(source string, target []string) bool {
				return !slices.Contains(target, source)
			},
			messageTemplate: "actual value %s, allowed (%s): %w",
			name:            "in",
			baseError:       ErrValueIsNotAllowed,
		},
	}
}

func NewInIntConstraint() core.ConstraintInterface {
	return &comparableConstraint[int64]{
		baseConstraint[int64, []int64]{
			kindTypes: []reflect.Kind{
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			},
			prepareValue: func(value reflect.Value) int64 {
				return value.Int()
			},
			prepareParam: func(params string) ([]int64, error) {
				pieces := strings.Split(params, ",")
				result := make([]int64, len(pieces))
				var parsErr error
				for i, piece := range pieces {
					result[i], parsErr = strconv.ParseInt(piece, 10, 64)
					if parsErr != nil {
						return nil, parsErr
					}
				}

				return result, nil
			},
			compareValues: func(source int64, target []int64) bool {
				return !slices.Contains(target, source)
			},
			messageTemplate: "actual value %s, allowed (%s): %w",
			name:            "in",
			baseError:       ErrValueIsNotAllowed,
		},
	}
}

func NewInFloatConstraint() core.ConstraintInterface {
	return &comparableConstraint[float64]{
		baseConstraint[float64, []float64]{
			kindTypes: []reflect.Kind{
				reflect.Float32, reflect.Float64,
			},
			prepareValue: func(value reflect.Value) float64 {
				return value.Float()
			},
			prepareParam: func(params string) ([]float64, error) {
				pieces := strings.Split(params, ",")
				result := make([]float64, len(pieces))
				var parsErr error
				for i, piece := range pieces {
					result[i], parsErr = strconv.ParseFloat(piece, 64)
					if parsErr != nil {
						return nil, parsErr
					}
				}

				return result, nil
			},
			compareValues: func(source float64, target []float64) bool {
				return !slices.Contains(target, source)
			},
			messageTemplate: "actual value %s, allowed (%s): %w",
			name:            "in",
			baseError:       ErrValueIsNotAllowed,
		},
	}
}

func init() {
	core.RegisterConstraint(NewInStringConstraint())
	core.RegisterConstraint(NewInIntConstraint())
	core.RegisterConstraint(NewInFloatConstraint())
}
