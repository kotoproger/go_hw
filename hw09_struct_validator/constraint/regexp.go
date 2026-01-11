package constraint

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/kotoproger/go_hw/hw09structvalidator"
)

var ErrDoesNotMatchPattern = errors.New("does not match the pattern")

type RegexpConstraint struct {
	baseConstraint[string, *regexp.Regexp]
}

func NewRegexpConstraint() hw09structvalidator.ConstraintInterface {
	return RegexpConstraint{
		baseConstraint[string, *regexp.Regexp]{
			kindTypes: []reflect.Kind{reflect.String},
			prepareValue: func(value reflect.Value) string {
				return value.String()
			},
			prepareParam: regexp.Compile,
			compareValues: func(source string, target *regexp.Regexp) bool {
				return !target.MatchString(source)
			},
			messageTemplate: "value %s, pattern %s: %w",
			name:            "regexp",
			subname:         "string",
			baseError:       ErrDoesNotMatchPattern,
		},
	}
}

func init() {
	hw09structvalidator.RegisterConstraint(NewRegexpConstraint())
}
