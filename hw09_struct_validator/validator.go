package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kotoproger/go_hw/hw09structvalidator/core"
)

type ValidationErrors []core.ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, err := range v {
		builder.WriteString(
			err.Error() + "\n",
		)
	}

	return builder.String()
}

func Validate(value interface{}) error {
	errors := ValidationErrors{}
	valueReflection := reflect.ValueOf(value)
	if valueReflection.Kind() == reflect.Pointer {
		valueReflection = valueReflection.Elem()
	}

	if valueReflection.Kind() != reflect.Struct {
		return nil
	}

	typeReflection := valueReflection.Type()

	for fieldNumber := 0; fieldNumber < typeReflection.NumField(); fieldNumber++ {
		tagValue, isset := typeReflection.Field(fieldNumber).Tag.Lookup("validate")
		if !isset {
			continue
		}

		validationParams := strings.Split(tagValue, "|")
		for _, validationParam := range validationParams {
			ConstraintInfo := strings.Split(validationParam, ":")
			if len(ConstraintInfo) != 2 {

				return core.ErrUnsupportedConstraintParams(fmt.Errorf("invalid validation params `%s`", validationParam))
			}
			constraintTypes, ok := core.Constraints[ConstraintInfo[0]]
			if !ok {

				return core.ErrUnknownConstraint(fmt.Errorf("invalid validation params `%s`", validationParam))
			}
			typeConstraint, ok := constraintTypes[typeReflection.Field(fieldNumber).Type.Kind()]
			if !ok {

				return core.ErrUnsupportedValueType
			}

			constraintErrors, validationError := typeConstraint.Validate(
				valueReflection.Field(fieldNumber),
				ConstraintInfo[1],
			)

			if validationError != nil {
				return validationError
			}

			for _, constraintError := range constraintErrors {
				errors = append(errors, core.ValidationError{
					Field: typeReflection.Field(fieldNumber).Name,
					Err:   constraintError,
				})
			}

		}
	}
	return errors
}
