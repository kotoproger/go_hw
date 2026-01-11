package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/kotoproger/go_hw/hw09structvalidator/constraint"
	"github.com/kotoproger/go_hw/hw09structvalidator/core"
)

type ValidationErrors []core.ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, err := range v {
		builder.WriteString(
			fmt.Sprintf("Validation error for field %s: %s\n", err.Field, err.Err))
	}

	return builder.String()
}

func Validate(value interface{}) ValidationErrors {
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
				errors = append(errors, core.ValidationError{
					Field: typeReflection.Name(),
					Err:   core.ErrUnsupportedConstraintParams(fmt.Errorf("invalid validation params `%s`", validationParam)),
				})

				continue
			}
			constraintTypes, ok := core.Constraints[ConstraintInfo[0]]
			if !ok {
				errors = append(errors, core.ValidationError{
					Field: typeReflection.Name(),
					Err:   core.ErrUnknownConstraint(fmt.Errorf("invalid validation params `%s`", validationParam)),
				})

				continue
			}
			typeConstraint, ok := constraintTypes[typeReflection.Field(fieldNumber).Type.Kind()]
			if !ok {
				errors = append(errors, core.ValidationError{
					Field: typeReflection.Name(),
					Err:   core.ErrUnsupportedValueType,
				})

				continue
			}

			constraintErrors, validationError := typeConstraint.Validate(
				valueReflection.Field(fieldNumber),
				ConstraintInfo[1],
			)

			for _, constraintError := range constraintErrors {
				errors = append(errors, core.ValidationError{
					Field: typeReflection.Name(),
					Err:   constraintError,
				})
			}
			if validationError != nil {
				errors = append(errors, core.ValidationError{
					Field: typeReflection.Name(),
					Err:   validationError,
				})
			}
		}
	}
	return errors
}
