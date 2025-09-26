package validation

import (
	"donbarrigon/new/internal/utils/fm"
	"fmt"
	"reflect"
	"strconv"
)

// In valida que el valor esté presente en la lista de valores permitidos
func In(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 1 {
		return "Se requiere al menos un valor permitido", fm.Placeholder{}, true
	}

	ph := fm.Placeholder{"values": fmt.Sprintf("%v", params)}

	switch value.Kind() {
	case reflect.String:
		valueStr := value.String()
		for _, param := range params {
			if valueStr == param {
				return "", nil, false
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueInt := value.Int()
		for _, param := range params {
			if paramInt, e := strconv.ParseInt(param, 10, 64); e == nil {
				if valueInt == paramInt {
					return "", nil, false
				}
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valueUint := value.Uint()
		for _, param := range params {
			if paramUint, e := strconv.ParseUint(param, 10, 64); e == nil {
				if valueUint == paramUint {
					return "", nil, false
				}
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	case reflect.Float32, reflect.Float64:
		valueFloat := value.Float()
		for _, param := range params {
			if paramFloat, e := strconv.ParseFloat(param, 64); e == nil {
				if valueFloat == paramFloat {
					return "", nil, false
				}
			}
		}
		return "El campo :field debe ser uno de los valores :values", ph, true

	case reflect.Slice, reflect.Array:
		// Verificar que todos los elementos del array/slice estén en la lista permitida
		for i := 0; i < value.Len(); i++ {
			element := value.Index(i)
			elementFound := false

			switch element.Kind() {
			case reflect.String:
				elementStr := element.String()
				for _, param := range params {
					if elementStr == param {
						elementFound = true
						break
					}
				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				elementInt := element.Int()
				for _, param := range params {
					if paramInt, e := strconv.ParseInt(param, 10, 64); e == nil {
						if elementInt == paramInt {
							elementFound = true
							break
						}
					}
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				elementUint := element.Uint()
				for _, param := range params {
					if paramUint, e := strconv.ParseUint(param, 10, 64); e == nil {
						if elementUint == paramUint {
							elementFound = true
							break
						}
					}
				}

			case reflect.Float32, reflect.Float64:
				elementFloat := element.Float()
				for _, param := range params {
					if paramFloat, e := strconv.ParseFloat(param, 64); e == nil {
						if elementFloat == paramFloat {
							elementFound = true
							break
						}
					}
				}

			default:
				return "Tipo de elemento no soportado en el array/slice para validación In", ph, true
			}

			if !elementFound {
				return "Todos los elementos del campo :field deben ser valores permitidos: :values", ph, true
			}
		}
		return "", nil, false

	default:
		return "Tipo no soportado para validación In", ph, true
	}
}

// Nin valida que el valor NO esté presente en la lista de valores prohibidos
func Nin(value reflect.Value, params ...string) (string, fm.Placeholder, bool) {
	if len(params) < 1 {
		return "Se requiere al menos un valor prohibido", fm.Placeholder{}, true
	}

	ph := fm.Placeholder{"values": fmt.Sprintf("%v", params)}

	switch value.Kind() {
	case reflect.String:
		valueStr := value.String()
		for _, param := range params {
			if valueStr == param {
				return "El campo :field no puede ser uno de los valores :values", ph, true
			}
		}
		return "", nil, false

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valueInt := value.Int()
		for _, param := range params {
			if paramInt, e := strconv.ParseInt(param, 10, 64); e == nil {
				if valueInt == paramInt {
					return "El campo :field no puede ser uno de los valores :values", ph, true
				}
			}
		}
		return "", nil, false

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valueUint := value.Uint()
		for _, param := range params {
			if paramUint, e := strconv.ParseUint(param, 10, 64); e == nil {
				if valueUint == paramUint {
					return "El campo :field no puede ser uno de los valores :values", ph, true
				}
			}
		}
		return "", nil, false

	case reflect.Float32, reflect.Float64:
		valueFloat := value.Float()
		for _, param := range params {
			if paramFloat, e := strconv.ParseFloat(param, 64); e == nil {
				if valueFloat == paramFloat {
					return "El campo :field no puede ser uno de los valores :values", ph, true
				}
			}
		}
		return "", nil, false

	case reflect.Slice, reflect.Array:
		// Verificar que ningún elemento del array/slice esté en la lista prohibida
		for i := 0; i < value.Len(); i++ {
			element := value.Index(i)

			switch element.Kind() {
			case reflect.String:
				elementStr := element.String()
				for _, param := range params {
					if elementStr == param {
						return "Ningún elemento del campo :field puede ser uno de los valores :values", ph, true
					}
				}

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				elementInt := element.Int()
				for _, param := range params {
					if paramInt, e := strconv.ParseInt(param, 10, 64); e == nil {
						if elementInt == paramInt {
							return "Ningún elemento del campo :field puede ser uno de los valores :values", ph, true
						}
					}
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				elementUint := element.Uint()
				for _, param := range params {
					if paramUint, e := strconv.ParseUint(param, 10, 64); e == nil {
						if elementUint == paramUint {
							return "Ningún elemento del campo :field puede ser uno de los valores :values", ph, true
						}
					}
				}

			case reflect.Float32, reflect.Float64:
				elementFloat := element.Float()
				for _, param := range params {
					if paramFloat, e := strconv.ParseFloat(param, 64); e == nil {
						if elementFloat == paramFloat {
							return "Ningún elemento del campo :field puede ser uno de los valores :values", ph, true
						}
					}
				}

			default:
				return "Tipo de elemento no soportado en el array/slice para validación Nin", ph, true
			}
		}
		return "", nil, false

	default:
		return "Tipo no soportado para validación Nin", ph, true
	}
}
