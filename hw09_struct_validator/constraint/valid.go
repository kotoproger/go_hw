package constraint

// import (
//	"reflect"
//
//	"github.com/kotoproger/go_hw/hw09structvalidator"
//)
//
// type validConstraint struct {
//	baseConstraint[any, any]
//}
//
// func NewValidConstraint() hw09structvalidator.ConstraintInterface {
//	return validConstraint{
//		baseConstraint[any, any]{
//			kindTypes: []reflect.Kind{reflect.Struct},
//			prepareValue: func(value reflect.Value) any {
//				return value.Interface()
//			},
//			prepareParam: func(_ string) (any, error) {
//				return nil, nil
//			},
//			compareValues: func(source any, _ any) bool {
//				hw09structvalidator.Validate(source)
//			},
//		},
//	}
//}
