package conversion

import (
	"fmt"
	"reflect"
)

// Verifies whether a conversion function has a correct signature.
func verifyConversionFunctionSignature(ft reflect.Type) error {
	if ft.Kind() != reflect.Func {
		return fmt.Errorf("expected func, got: %v", ft)
	}
	if ft.NumIn() != 3 {
		return fmt.Errorf("expected three 'in' params, got: %v", ft)
	}
	if ft.NumOut() != 1 {
		return fmt.Errorf("expected one 'out' param, got: %v", ft)
	}
	if ft.In(0).Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer arg for 'in' param 0, got: %v", ft)
	}
	if ft.In(1).Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer arg for 'in' param 1, got: %v", ft)
	}
	scopeType := Scope(nil)
	if e, a := reflect.TypeOf(&scopeType).Elem(), ft.In(2); e != a {
		return fmt.Errorf("expected '%v' arg for 'in' param 2, got '%v' (%v)", e, a, ft)
	}
	var forErrorType error
	// This convolution is necessary, otherwise TypeOf picks up on the fact
	// that forErrorType is nil.
	errorType := reflect.TypeOf(&forErrorType).Elem()
	if ft.Out(0) != errorType {
		return fmt.Errorf("expected error return, got: %v", ft)
	}
	return nil
}

// reflectFields lists all struct fields
func reflectFields(rt reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			v := rt.Field(i)
			if v.Anonymous {
				fields = append(fields, reflectFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}
	return fields
}

func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}
