package validate

import (
	"donbarrigon/new/internal/utils/fm"
	"reflect"
)

// Required valida que el valor no sea falsy o un valor cero
func Required(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	ph := fm.Placeholder{}

	if !value.IsValid() {
		return "El campo :field es requerido", ph, true
	}

	switch value.Kind() {
	case reflect.String:
		if value.String() == "" {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == 0 {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value.Uint() == 0 {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Float32, reflect.Float64:
		if value.Float() == 0.0 {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Bool:
		if !value.Bool() {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Slice, reflect.Array:
		if value.IsNil() || value.Len() == 0 {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Map:
		if value.IsNil() || value.Len() == 0 {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Ptr, reflect.Interface:
		if value.IsNil() {
			return "El campo :field es requerido", ph, true
		}
		return Required(value.Elem(), params...)

	case reflect.Chan, reflect.Func:
		if value.IsNil() {
			return "El campo :field es requerido", ph, true
		}

	case reflect.Struct:
		zeroValue := reflect.Zero(value.Type())
		if reflect.DeepEqual(value.Interface(), zeroValue.Interface()) {
			return "El campo :field es requerido", ph, true
		}

	default:
		if value.IsZero() {
			return "El campo :field es requerido", ph, true
		}
	}

	return "", nil, false
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZeroValue(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		// Para time.Time u otros structs
		return v.IsZero()
	default:
		return false
	}
}
